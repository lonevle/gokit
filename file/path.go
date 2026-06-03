package file

import (
	"os"
	"path/filepath"
)

// GetProgramPath 获取程序所在目录
// 返回:
//   - 程序所在目录路径
//   - 错误信息（如获取失败）
func GetProgramPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

// GetProgramName 获取程序名称
// 返回:
//   - 程序文件名（不含路径）
//   - 错误信息（如获取失败）
func GetProgramName() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Base(exePath), nil
}
