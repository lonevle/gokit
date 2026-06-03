package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"os"
)

// GCM 格式: [4:nonceLen][变长:nonce][变长:ciphertext+tag]

// GCM_Encrypt 加密字节切片
func GCM_Encrypt(plaintext, key []byte) ([]byte, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := randomBytes(gcm.NonceSize())
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// 输出: nonceLen(4 bytes) || nonce || ciphertext
	out := make([]byte, 4+len(nonce)+len(ciphertext))
	binary.BigEndian.PutUint32(out[:4], uint32(len(nonce)))
	copy(out[4:], nonce)
	copy(out[4+len(nonce):], ciphertext)

	return out, nil
}

// GCM_Decrypt 解密字节切片
func GCM_Decrypt(data, key []byte) ([]byte, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	if len(data) < 4 {
		return nil, ErrShortData
	}

	nonceSize := binary.BigEndian.Uint32(data[:4])
	if len(data) < int(4+nonceSize) {
		return nil, ErrShortData
	}

	nonce := data[4 : 4+nonceSize]
	ciphertext := data[4+nonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, int(nonceSize))
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
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
