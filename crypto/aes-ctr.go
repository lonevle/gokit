package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"os"
)

// CTR 文件头: [16:IV]
const ctrIVLen = aes.BlockSize // 16
// CTRHMAC 格式: [16:IV][32:HMAC][变长:ciphertext]
const ctrhmacHMACLen = sha256.Size // 32

// deriveTwoKeys 从单个密钥派生两个独立密钥
func deriveTwoKeys(key []byte) (encKey, macKey []byte) {
	if len(key) >= 32 {
		mid := len(key) / 2
		return key[:mid], key[mid:]
	}
	// 短密钥用 HMAC 派生
	h := hmac.New(sha256.New, key)
	h.Write([]byte("enc"))
	encKey = h.Sum(nil)[:len(key)]

	h.Reset()
	h.Write([]byte("mac"))
	macKey = h.Sum(nil)

	return encKey, macKey
}

// CTR_Encrypt 加密字节切片
func CTR_Encrypt(plaintext, key []byte) ([]byte, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}

	iv, err := randomBytes(ctrIVLen)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	// 输出: IV + ciphertext
	return append(iv, ciphertext...), nil
}

// CTR_Decrypt 解密字节切片
func CTR_Decrypt(data, key []byte) ([]byte, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	if len(data) < ctrIVLen {
		return nil, ErrShortData
	}

	iv, ciphertext := data[:ctrIVLen], data[ctrIVLen:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// CTR_EncryptFile 加密文件（流式，低内存）
func CTR_EncryptFile(inPath, outPath string, key []byte) error {
	if err := checkKey(key); err != nil {
		return err
	}

	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	tmpPath := outPath + ".tmp"
	outFile, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	cleanup := true
	defer func() {
		outFile.Close()
		if cleanup {
			os.Remove(tmpPath)
		}
	}()

	iv, err := randomBytes(ctrIVLen)
	if err != nil {
		return err
	}
	if _, err := outFile.Write(iv); err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)

	if _, err = io.Copy(outFile, &cipher.StreamReader{S: stream, R: inFile}); err != nil {
		return err
	}
	if err := outFile.Close(); err != nil {
		return err
	}

	cleanup = false
	return os.Rename(tmpPath, outPath)
}

// CTR_DecryptFile 解密文件（流式，低内存）
func CTR_DecryptFile(inPath, outPath string, key []byte) error {
	if err := checkKey(key); err != nil {
		return err
	}

	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	tmpPath := outPath + ".tmp"
	outFile, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	cleanup := true
	defer func() {
		outFile.Close()
		if cleanup {
			os.Remove(tmpPath)
		}
	}()

	iv := make([]byte, ctrIVLen)
	if _, err := io.ReadFull(inFile, iv); err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)

	if _, err = io.Copy(outFile, &cipher.StreamReader{S: stream, R: inFile}); err != nil {
		return err
	}
	if err := outFile.Close(); err != nil {
		return err
	}

	cleanup = false
	return os.Rename(tmpPath, outPath)
}

// CTRHMAC_Encrypt 加密并认证字节切片
func CTRHMAC_Encrypt(plaintext, key []byte) ([]byte, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}

	encKey, macKey := deriveTwoKeys(key)

	iv, err := randomBytes(ctrIVLen)
	if err != nil {
		return nil, err
	}

	// 加密
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	// HMAC(IV || ciphertext)
	h := hmac.New(sha256.New, macKey)
	h.Write(iv)
	h.Write(ciphertext)
	tag := h.Sum(nil)

	// 输出: IV || HMAC || ciphertext
	out := make([]byte, 0, ctrIVLen+ctrhmacHMACLen+len(ciphertext))
	out = append(out, iv...)
	out = append(out, tag...)
	out = append(out, ciphertext...)

	return out, nil
}

// CTRHMAC_Decrypt 验证并解密字节切片
func CTRHMAC_Decrypt(data, key []byte) ([]byte, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	if len(data) < ctrIVLen+ctrhmacHMACLen {
		return nil, ErrShortData
	}

	encKey, macKey := deriveTwoKeys(key)

	iv := data[:ctrIVLen]
	storedTag := data[ctrIVLen : ctrIVLen+ctrhmacHMACLen]
	ciphertext := data[ctrIVLen+ctrhmacHMACLen:]

	// 验证 HMAC
	h := hmac.New(sha256.New, macKey)
	h.Write(iv)
	h.Write(ciphertext)
	if !hmac.Equal(h.Sum(nil), storedTag) {
		return nil, ErrCorrupted
	}

	// 解密
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// CTRHMAC_EncryptFile 加密并认证文件（流式，低内存）
func CTRHMAC_EncryptFile(inPath, outPath string, key []byte) error {
	if err := checkKey(key); err != nil {
		return err
	}

	encKey, macKey := deriveTwoKeys(key)

	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	tmpPath := outPath + ".tmp"
	outFile, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	cleanup := true
	defer func() {
		outFile.Close()
		if cleanup {
			os.Remove(tmpPath)
		}
	}()

	iv, err := randomBytes(ctrIVLen)
	if err != nil {
		return err
	}
	if _, err := outFile.Write(iv); err != nil {
		return err
	}

	// 预留 HMAC 位置
	hmacPos, _ := outFile.Seek(0, io.SeekCurrent)
	if _, err := outFile.Write(make([]byte, ctrhmacHMACLen)); err != nil {
		return err
	}

	// 加密并计算 HMAC
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)

	h := hmac.New(sha256.New, macKey)
	h.Write(iv)

	// 流式处理
	buf := make([]byte, 1024*1024)
	reader := &cipher.StreamReader{S: stream, R: inFile}

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			h.Write(buf[:n])
			if _, werr := outFile.Write(buf[:n]); werr != nil {
				return werr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// 回填 HMAC
	tag := h.Sum(nil)
	if _, err := outFile.Seek(hmacPos, io.SeekStart); err != nil {
		return err
	}
	if _, err := outFile.Write(tag); err != nil {
		return err
	}

	if err := outFile.Close(); err != nil {
		return err
	}

	cleanup = false
	return os.Rename(tmpPath, outPath)
}

// CTRHMAC_DecryptFile 验证并解密文件（流式，低内存）
func CTRHMAC_DecryptFile(inPath, outPath string, key []byte) error {
	if err := checkKey(key); err != nil {
		return err
	}

	encKey, macKey := deriveTwoKeys(key)

	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	iv := make([]byte, ctrIVLen)
	if _, err := io.ReadFull(inFile, iv); err != nil {
		return err
	}

	storedTag := make([]byte, ctrhmacHMACLen)
	if _, err := io.ReadFull(inFile, storedTag); err != nil {
		return err
	}

	// 先验证 HMAC（流式读取，不写入）
	h := hmac.New(sha256.New, macKey)
	h.Write(iv)

	buf := make([]byte, 1024*1024)
	for {
		n, err := inFile.Read(buf)
		if n > 0 {
			h.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	if !hmac.Equal(h.Sum(nil), storedTag) {
		return ErrCorrupted
	}

	// 验证通过，复用句柄解密
	if _, err := inFile.Seek(int64(ctrIVLen+ctrhmacHMACLen), io.SeekStart); err != nil {
		return err
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return err
	}
	stream := cipher.NewCTR(block, iv)

	_, err = io.Copy(outFile, &cipher.StreamReader{S: stream, R: inFile})
	return err
}
