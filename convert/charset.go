// Package convert 提供编码识别与统一转码（UTF-8 / GBK / UTF-16 / UTF-32）。
//
// 主要解决 Windows 环境的中文乱码问题：
//   - cmd / powershell 控制台输出默认 GBK（cp936）
//   - 部分历史配置文件为 GBK 或带 BOM 的 UTF-8 / UTF-16 编码
//
// 本包只做纯转换，不进行任何文件 IO。
package convert

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// ToUTF8 识别字节流编码并统一转为 UTF-8。
//
// 检测顺序：
//  1. UTF-32 BOM → 按字节序解码（4 字节 BOM 先于 2 字节判断，避免 FF FE 00 00 被当成 UTF-16）
//  2. UTF-8 BOM → 剥离 BOM 头
//  3. UTF-16 BOM → 按字节序解码（须在 GBK 之前判断：FF FE 也是合法 GBK 序列，会被误解码成乱码）
//  4. 合法 UTF-8 → 原样返回（纯 ASCII 也在此命中，零拷贝）
//  5. 含 NUL 字节 → 判定为二进制，返回错误（文本编码的 BOM 场景已在前面处理完）
//  6. GBK 解码成功 → 转为 UTF-8
//  7. 以上皆失败 → 返回错误
//
// 已知局限：无 BOM 的 GBK 内容若恰好构成合法 UTF-8，会被误判为 UTF-8（启发式固有缺陷）。
//
// 大文件建议使用流式版本 NewUTF8Reader，避免一次性加载到内存。
func ToUTF8(data []byte) ([]byte, error) {
	// UTF-32 大端序 BOM: 00 00 FE FF
	if isUTF32BigEndianBOM(data) {
		out, _, err := transform.Bytes(&utf32Decoder{littleEndian: false}, data)
		return out, err
	}

	// UTF-32 小端序 BOM: FF FE 00 00
	if isUTF32LittleEndianBOM(data) {
		out, _, err := transform.Bytes(&utf32Decoder{littleEndian: true}, data)
		return out, err
	}

	// UTF-8 BOM: EF BB BF，剥离后继续后续判断
	if isUTF8BOM(data) {
		data = data[3:]
	}

	// UTF-16 大端序 BOM: FE FF
	if isUTF16BigEndianBOM(data) {
		return unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder().Bytes(data)
	}

	// UTF-16 小端序 BOM: FF FE
	if isUTF16LittleEndianBOM(data) {
		return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder().Bytes(data)
	}

	// 已是合法 UTF-8（含纯 ASCII），原样返回
	if utf8.Valid(data) {
		return data, nil
	}

	// 非 UTF-8 且含 NUL 字节，基本可断定为二进制（GBK 解码过于宽容，
	// 仅靠解码失败无法识别二进制，参考 file/git 的二进制判定方式）
	if bytes.IndexByte(data, 0x00) >= 0 {
		return nil, errors.New("convert: 疑似二进制数据（含 NUL 字节）")
	}

	// 尝试按 GBK 解码（Windows 控制台默认 cp936）
	if out, err := simplifiedchinese.GBK.NewDecoder().Bytes(data); err == nil {
		return out, nil
	}

	return nil, errors.New("convert: 无法识别的编码（非 UTF-8/GBK/UTF-16/UTF-32）")
}

// NewUTF8Reader 包装任意编码的读取流，返回转为 UTF-8 的读取流（流式版本）。
//
// 内存占用为 O(缓冲区)，适合处理大文件：
//
//	f, _ := os.Open("gbk.log")
//	defer f.Close()
//	r, err := convert.NewUTF8Reader(f)
//	if err != nil {
//		// 处理错误
//	}
//	io.Copy(os.Stdout, r) // 或用 bufio.Scanner 逐行处理
//
// 处理流程：
//   - 头部带 BOM（UTF-32/UTF-16/UTF-8）→ 按对应编码解码
//   - 无 BOM → 采样头部数据判断：合法 UTF-8 直接透传；含 NUL 字节判定二进制返回错误；否则按 GBK 解码
//
// 注意：无 BOM 时的采样判断是启发式，准确率略低于 ToUTF8 的全量判断；
// 流式解码错误会在后续 Read 时返回。
func NewUTF8Reader(r io.Reader) (io.Reader, error) {
	br := bufio.NewReader(r)

	// Peek 失败只可能是数据不足（EOF），按实际读到的长度判断即可
	head, err := br.Peek(4)
	if err != nil && err != io.EOF {
		return nil, err
	}

	switch {
	case isUTF32BigEndianBOM(head):
		return transform.NewReader(br, &utf32Decoder{littleEndian: false}), nil
	case isUTF32LittleEndianBOM(head):
		return transform.NewReader(br, &utf32Decoder{littleEndian: true}), nil
	case isUTF8BOM(head):
		// 剥离 BOM 后剩余内容已是 UTF-8，直接透传
		br.Discard(3)
		return br, nil
	case isUTF16BigEndianBOM(head):
		return transform.NewReader(br, unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder()), nil
	case isUTF16LittleEndianBOM(head):
		return transform.NewReader(br, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()), nil
	}

	// 无 BOM：采样头部判断 UTF-8 / GBK
	sample, err := br.Peek(1024)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if isLikelyUTF8(sample) {
		return br, nil
	}
	// 采样数据非 UTF-8 且含 NUL 字节，基本可断定为二进制
	if bytes.IndexByte(sample, 0x00) >= 0 {
		return nil, errors.New("convert: 疑似二进制数据（含 NUL 字节）")
	}
	return transform.NewReader(br, simplifiedchinese.GBK.NewDecoder()), nil
}

// utf32Decoder 实现 transform.Transformer，将 UTF-32 流解码为 UTF-8 流。
// x/text 未提供 UTF-32 实现，这里手动实现，ToUTF8 与 NewUTF8Reader 共用同一套逻辑。
type utf32Decoder struct {
	littleEndian bool // 字节序：true 小端，false 大端
	bomSkipped   bool // 是否已处理过开头首个码点（用于跳过 BOM）
}

// Reset 实现 transform.Transformer 接口，复用前重置状态
func (d *utf32Decoder) Reset() {
	d.bomSkipped = false
}

// Transform 实现 transform.Transformer 接口，逐 4 字节解码一个码点
func (d *utf32Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for nSrc < len(src) {
		// 剩余不足 4 字节：等待更多数据；流已结束则数据非法
		if len(src)-nSrc < 4 {
			if atEOF {
				return nDst, nSrc, errors.New("convert: 非法的 UTF-32 数据（末尾不足 4 字节）")
			}
			return nDst, nSrc, transform.ErrShortSrc
		}

		var v uint32
		if d.littleEndian {
			v = uint32(src[nSrc]) | uint32(src[nSrc+1])<<8 | uint32(src[nSrc+2])<<16 | uint32(src[nSrc+3])<<24
		} else {
			v = uint32(src[nSrc])<<24 | uint32(src[nSrc+1])<<16 | uint32(src[nSrc+2])<<8 | uint32(src[nSrc+3])
		}
		nSrc += 4

		// 首个码点为 BOM（0xFEFF）时跳过
		if !d.bomSkipped {
			d.bomSkipped = true
			if v == 0xFEFF {
				continue
			}
		}

		// 码点超过 Unicode 上限或落在代理区，均为非法
		if v > utf8.MaxRune || (v >= 0xD800 && v <= 0xDFFF) {
			return nDst, nSrc, fmt.Errorf("convert: 非法的 UTF-32 码点 0x%X", v)
		}

		// 目标缓冲区不足时回退本码点并归还控制权（单个码点转 UTF-8 最多 4 字节）
		if len(dst)-nDst < 4 {
			nSrc -= 4
			return nDst, nSrc, transform.ErrShortDst
		}
		nDst += utf8.EncodeRune(dst[nDst:], rune(v))
	}
	return nDst, nSrc, nil
}

// isLikelyUTF8 判断采样数据是否为合法 UTF-8。
// 末尾被采样截断的多字节序列不视为非法：只要截断点之前的部分完整合法即认定通过。
func isLikelyUTF8(sample []byte) bool {
	for i := 0; i < len(sample); {
		r, size := utf8.DecodeRune(sample[i:])
		// 遇到非法字节：若只是序列被采样截断（不完整）则放行，否则不是 UTF-8
		if r == utf8.RuneError && size == 1 {
			return !utf8.FullRune(sample[i:])
		}
		i += size
	}
	return true
}

// isUTF32BigEndianBOM 判断是否为 UTF-32 大端序 BOM (00 00 FE FF)
func isUTF32BigEndianBOM(b []byte) bool {
	return len(b) >= 4 && b[0] == 0x00 && b[1] == 0x00 && b[2] == 0xFE && b[3] == 0xFF
}

// isUTF32LittleEndianBOM 判断是否为 UTF-32 小端序 BOM (FF FE 00 00)
func isUTF32LittleEndianBOM(b []byte) bool {
	return len(b) >= 4 && b[0] == 0xFF && b[1] == 0xFE && b[2] == 0x00 && b[3] == 0x00
}

// isUTF8BOM 判断是否为 UTF-8 BOM (EF BB BF)
func isUTF8BOM(b []byte) bool {
	return len(b) >= 3 && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF
}

// isUTF16BigEndianBOM 判断是否为 UTF-16 大端序 BOM (FE FF)
func isUTF16BigEndianBOM(b []byte) bool {
	return len(b) >= 2 && b[0] == 0xFE && b[1] == 0xFF
}

// isUTF16LittleEndianBOM 判断是否为 UTF-16 小端序 BOM (FF FE)
func isUTF16LittleEndianBOM(b []byte) bool {
	return len(b) >= 2 && b[0] == 0xFF && b[1] == 0xFE
}
