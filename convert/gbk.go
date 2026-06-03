package convert

import (
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// GBKToUTF8 将 GBK 编码的字节数组转换为 UTF-8 编码
// 如果输入已经是 UTF-8 编码，则直接返回
// 参数:
//   - data: 输入的字节数组
//
// 返回:
//   - UTF-8 编码的字节数组
//   - 转换错误
//
// 大文件建议使用流式转换，避免一次性加载到内存中
// f, _ := os.Open("gbk")
// defer f.Close()
// reader := transform.NewReader(f, simplifiedchinese.GBK.NewDecoder())
// scanner := bufio.NewScanner(reader)
//
//	for scanner.Scan() {
//	    // 处理一行，内存只保留一行
//	}
func GBKToUTF8(data []byte) ([]byte, error) {
	// 如果已经是 UTF-8 编码，直接返回
	if utf8.Valid(data) {
		return data, nil
	}

	// 使用 GBK 解码器转换
	return simplifiedchinese.GBK.NewDecoder().Bytes(data)
}
