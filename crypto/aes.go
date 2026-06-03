package crypto

import (
	"crypto/rand"
	"fmt"
)

var (
	ErrInvalidKey = fmt.Errorf("invalid key size, must be 16, 24, or 32 bytes")
	ErrCorrupted  = fmt.Errorf("data corrupted, tampered, or wrong key")
	ErrShortData  = fmt.Errorf("data too short")
)

// checkKey 验证 AES 密钥长度是否合法
func checkKey(key []byte) error {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return ErrInvalidKey
	}
	return nil
}

// randomBytes 生成指定长度的随机字节数组
func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
