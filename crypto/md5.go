package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// MD5 计算字符串的 MD5 哈希值
// 参数:
//   - str: 待计算的字符串
//
// 返回:
//   - 32位小写的 MD5 哈希字符串
func MD5(str string) string {
	sum := md5.Sum([]byte(str))
	return hex.EncodeToString(sum[:])
}

// FileMD5 计算文件的 MD5 哈希值
// 参数:
//   - path: 文件路径
//
// 返回:
//   - 32位小写的 MD5 哈希字符串
//   - 错误信息（如文件不存在或读取失败）
func FileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
