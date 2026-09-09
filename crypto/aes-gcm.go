package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"
)

// GCM 报文格式: [1B 版本][12B nonce][NB 密文+tag]
const (
	gcmVersionV1 = 1  // 报文格式版本号
	gcmNonceLen  = 12 // GCM 标准 nonce 长度
)

// newGCM 根据密钥创建 AES-GCM 实例
func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// GCM_Encrypt 加密字节切片
//
// GCM 为认证加密模式，自带完整性保护，密文被篡改时解密会失败。
// 输出格式: [1B 版本][12B nonce][NB 密文+tag]
func GCM_Encrypt(plaintext, key []byte) ([]byte, error) {
	return gcmEncrypt(plaintext, nil, key)
}

// GCM_EncryptWithAAD 加密字节切片并绑定附加认证数据
//
// aad 参与完整性校验但不加密，用于把密文之外的关联字段（如协议头、
// 用户 ID 等明文字段）纳入防篡改范围，这些字段被篡改时解密失败。
// 解密时必须传入与加密时一致的 aad。
func GCM_EncryptWithAAD(plaintext, aad, key []byte) ([]byte, error) {
	return gcmEncrypt(plaintext, aad, key)
}

// gcmEncrypt AES-GCM 加密内部实现
func gcmEncrypt(plaintext, aad, key []byte) ([]byte, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}

	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	// 预分配: 1(版本) + 12(nonce) + 明文 + tag，nonce 直接生成在输出缓冲区内
	out := make([]byte, 1+gcmNonceLen, 1+gcmNonceLen+len(plaintext)+gcm.Overhead())
	out[0] = gcmVersionV1
	nonce := out[1 : 1+gcmNonceLen]
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("crypto: 生成 nonce 失败: %w", err)
	}
	return gcm.Seal(out, nonce, plaintext, aad), nil
}

// GCM_Decrypt 解密字节切片，自动校验 tag
//
// 报文格式不匹配、版本不支持、密文被篡改或密钥错误均返回错误。
func GCM_Decrypt(data, key []byte) ([]byte, error) {
	return gcmDecrypt(data, nil, key)
}

// GCM_DecryptWithAAD 解密字节切片，aad 必须与加密时一致，否则解密失败
func GCM_DecryptWithAAD(data, aad, key []byte) ([]byte, error) {
	return gcmDecrypt(data, aad, key)
}

// gcmDecrypt AES-GCM 解密内部实现
func gcmDecrypt(data, aad, key []byte) ([]byte, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	if len(data) < 1+gcmNonceLen {
		return nil, ErrShortData
	}
	if data[0] != gcmVersionV1 {
		return nil, fmt.Errorf("crypto: 不支持的报文版本 %d", data[0])
	}

	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := data[1 : 1+gcmNonceLen]
	plaintext, err := gcm.Open(nil, nonce, data[1+gcmNonceLen:], aad)
	if err != nil {
		return nil, ErrCorrupted
	}
	return plaintext, nil
}

// GCM_EncryptFile 加密文件（一次性读取，适合中小文件）
func GCM_EncryptFile(inPath, outPath string, key []byte) error {
	if err := checkKey(key); err != nil {
		return err
	}

	plaintext, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}

	out, err := GCM_Encrypt(plaintext, key)
	if err != nil {
		return err
	}

	tmpPath := outPath + ".tmp"
	if err := os.WriteFile(tmpPath, out, 0600); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, outPath)
}

// GCM_DecryptFile 解密文件（一次性读取，适合中小文件）
func GCM_DecryptFile(inPath, outPath string, key []byte) error {
	if err := checkKey(key); err != nil {
		return err
	}

	data, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}

	plaintext, err := GCM_Decrypt(data, key)
	if err != nil {
		return err
	}

	tmpPath := outPath + ".tmp"
	if err := os.WriteFile(tmpPath, plaintext, 0600); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, outPath)
}
