package file

import (
	"os"

	"github.com/adhocore/jsonc"
	"github.com/lonevle/gokit/convert"
)

// Exists 检查路径是否存在
// 参数:
//   - path: 文件或目录路径
//
// 返回:
//   - true: 路径存在
//   - false: 路径不存在
//   - error: 其他错误（如权限不足）
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ReadFile 读取文件并统一转为 UTF-8（自动识别并处理 BOM 头和 GBK / UTF-16 等编码）
// 参数:
//   - filename: 文件路径
//
// 返回:
//   - UTF-8 编码的字节数组（已剥离 BOM）
//   - 错误信息
func ReadFile(filename string) ([]byte, error) {
	fileByte, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return convert.ToUTF8(fileByte)
}

// ReadJson 读取 JSON 文件，自动处理 BOM、GBK 编码转换并去除注释
// 参数:
//   - filename: JSON 文件路径
//
// 返回:
//   - 去除注释后的 UTF-8 编码字节数组
//   - 错误信息
func ReadJson(filename string) ([]byte, error) {
	u8b, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}

	jc := jsonc.New()
	return jc.Strip(u8b), nil
}
