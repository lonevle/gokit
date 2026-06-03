package convert

import (
	"os"
)

// SkipBOM 读取文件并跳过 BOM 头
// 支持 UTF-8、UTF-16 (大端/小端)、UTF-32 (大端/小端) 的 BOM 头
// 参数:
//   - filename: 文件路径
// 返回:
//   - 剔除 BOM 头后的字节数组
//   - 错误信息
func SkipBOM(filename string) ([]byte, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// UTF-32 大端序 BOM: 00 00 FE FF
	if len(file) >= 4 && isUTF32BigEndianBOM(file) {
		return file[4:], nil
	}

	// UTF-32 小端序 BOM: FF FE 00 00
	if len(file) >= 4 && isUTF32LittleEndianBOM(file) {
		return file[4:], nil
	}

	// UTF-8 BOM: EF BB BF
	if len(file) >= 3 && isUTF8BOM(file) {
		return file[3:], nil
	}

	// UTF-16 大端序 BOM: FE FF
	if len(file) >= 2 && isUTF16BigEndianBOM(file) {
		return file[2:], nil
	}

	// UTF-16 小端序 BOM: FF FE
	if len(file) >= 2 && isUTF16LittleEndianBOM(file) {
		return file[2:], nil
	}

	return file, nil
}

// isUTF32BigEndianBOM 判断是否为 UTF-32 大端序 BOM (00 00 FE FF)
func isUTF32BigEndianBOM(buf []byte) bool {
	return len(buf) >= 4 &&
		buf[0] == 0x00 && buf[1] == 0x00 &&
		buf[2] == 0xFE && buf[3] == 0xFF
}

// isUTF32LittleEndianBOM 判断是否为 UTF-32 小端序 BOM (FF FE 00 00)
func isUTF32LittleEndianBOM(buf []byte) bool {
	return len(buf) >= 4 &&
		buf[0] == 0xFF && buf[1] == 0xFE &&
		buf[2] == 0x00 && buf[3] == 0x00
}

// isUTF8BOM 判断是否为 UTF-8 BOM (EF BB BF)
func isUTF8BOM(buf []byte) bool {
	return len(buf) >= 3 &&
		buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF
}

// isUTF16BigEndianBOM 判断是否为 UTF-16 大端序 BOM (FE FF)
func isUTF16BigEndianBOM(buf []byte) bool {
	return len(buf) >= 2 && buf[0] == 0xFE && buf[1] == 0xFF
}

// isUTF16LittleEndianBOM 判断是否为 UTF-16 小端序 BOM (FF FE)
func isUTF16LittleEndianBOM(buf []byte) bool {
	return len(buf) >= 2 && buf[0] == 0xFF && buf[1] == 0xFE
}
