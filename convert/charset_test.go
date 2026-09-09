package convert

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// 各编码下 "你好" 的字节序列，供测试用例复用
var (
	utf16LE  = []byte{0xFF, 0xFE, 0x60, 0x4F, 0x7D, 0x59}                                 // UTF-16 小端序带 BOM
	utf16BE  = []byte{0xFE, 0xFF, 0x4F, 0x60, 0x59, 0x7D}                                 // UTF-16 大端序带 BOM
	utf32LE  = []byte{0xFF, 0xFE, 0x00, 0x00, 0x60, 0x4F, 0x00, 0x00, 0x7D, 0x59, 0x00, 0x00} // UTF-32 小端序带 BOM
	utf32BE  = []byte{0x00, 0x00, 0xFE, 0xFF, 0x00, 0x00, 0x4F, 0x60, 0x00, 0x00, 0x59, 0x7D} // UTF-32 大端序带 BOM
	gbk      = []byte{0xC4, 0xE3, 0xBA, 0xC3}                                            // GBK
	utf8BOM  = append([]byte{0xEF, 0xBB, 0xBF}, []byte("你好")...)                          // UTF-8 带 BOM
	longGBK  = bytes.Repeat(gbk, 1000)                                                   // 4000 字节，超过流式采样窗口
	longUTF8 = strings.Repeat("你好", 2000)                                                // 6000 字节，采样窗口必然截断多字节字符
)

// TestToUTF8 测试内存版编码识别与转换
func TestToUTF8(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
		want string
	}{
		{"纯ASCII", []byte("hello"), "hello"},
		{"UTF-8无BOM", []byte("你好"), "你好"},
		{"UTF-8带BOM", utf8BOM, "你好"},
		{"GBK", gbk, "你好"},
		{"长GBK", longGBK, strings.Repeat("你好", 1000)},
		{"UTF-16LE带BOM", utf16LE, "你好"},
		{"UTF-16BE带BOM", utf16BE, "你好"},
		{"UTF-32LE带BOM", utf32LE, "你好"},
		{"UTF-32BE带BOM", utf32BE, "你好"},
	}
	for _, c := range cases {
		got, err := ToUTF8(c.in)
		if err != nil {
			t.Errorf("%s: 意外错误: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s: 转换结果错误 (got %d 字节, want %d 字节)", c.name, len(got), len(c.want))
		}
	}

	// 二进制内容应返回错误
	if _, err := ToUTF8([]byte{0xFF, 0xC0, 0x00}); err == nil {
		t.Error("二进制输入应返回错误")
	}
}

// TestNewUTF8Reader 测试流式版编码识别与转换
func TestNewUTF8Reader(t *testing.T) {
	cases := []struct {
		name string
		in   []byte
		want string
	}{
		{"空流", nil, ""},
		{"GBK流", gbk, "你好"},
		{"长GBK流", longGBK, strings.Repeat("你好", 1000)},
		{"长UTF-8流", []byte(longUTF8), longUTF8},
		{"UTF-8BOM流", utf8BOM, "你好"},
		{"UTF-16LE流", utf16LE, "你好"},
		{"UTF-16BE流", utf16BE, "你好"},
		{"UTF-32LE流", utf32LE, "你好"},
		{"UTF-32BE流", utf32BE, "你好"},
	}
	for _, c := range cases {
		r, err := NewUTF8Reader(bytes.NewReader(c.in))
		if err != nil {
			t.Errorf("%s: 意外错误: %v", c.name, err)
			continue
		}
		got, err := io.ReadAll(r)
		if err != nil {
			t.Errorf("%s: 读取错误: %v", c.name, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf("%s: 转换结果错误 (got %d 字节, want %d 字节)", c.name, len(got), len(c.want))
		}
	}

	// UTF-32 流末尾不足 4 字节应在读取时报错
	r, _ := NewUTF8Reader(bytes.NewReader([]byte{0xFF, 0xFE, 0x00, 0x00, 0x60, 0x4F, 0x00}))
	if _, err := io.ReadAll(r); err == nil {
		t.Error("残缺的 UTF-32 流应返回错误")
	}

	// 二进制流应在创建时报错
	if _, err := NewUTF8Reader(bytes.NewReader([]byte{0xFF, 0xC0, 0x00, 0x01, 0x02, 0x03})); err == nil {
		t.Error("二进制流应返回错误")
	}
}
