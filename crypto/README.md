# crypto

加密工具包，提供 MD5 哈希计算、AES 加密/解密功能（支持 CTR、GCM、CTR+HMAC 模式）。

## 安装

```bash
go get github.com/lonevle/gokit/crypto
```

## 函数列表

### MD5

| 函数 | 签名 | 说明 |
|------|------|------|
| MD5 | `MD5(str string) string` | 计算字符串的 MD5 哈希值，返回 32 位小写字符串 |
| FileMD5 | `FileMD5(path string) (string, error)` | 计算文件的 MD5 哈希值，返回 32 位小写字符串 |

### AES-CTR（流式，低内存，无认证）

| 函数 | 签名 | 说明 |
|------|------|------|
| CTR_Encrypt | `CTR_Encrypt(plaintext, key []byte) ([]byte, error)` | 加密字节切片，输出格式: IV(16B) + ciphertext |
| CTR_Decrypt | `CTR_Decrypt(data, key []byte) ([]byte, error)` | 解密字节切片 |
| CTR_EncryptFile | `CTR_EncryptFile(inPath, outPath string, key []byte) error` | 流式加密文件，低内存占用 |
| CTR_DecryptFile | `CTR_DecryptFile(inPath, outPath string, key []byte) error` | 流式解密文件，低内存占用 |

### AES-GCM（一次性读取，自带认证，适合中小文件）

| 函数 | 签名 | 说明 |
|------|------|------|
| GCM_Encrypt | `GCM_Encrypt(plaintext, key []byte) ([]byte, error)` | 加密字节切片，输出格式: nonceLen(4B) + nonce + ciphertext+tag |
| GCM_Decrypt | `GCM_Decrypt(data, key []byte) ([]byte, error)` | 解密字节切片，自动验证 tag |
| GCM_EncryptFile | `GCM_EncryptFile(inPath, outPath string, key []byte) error` | 加密文件（一次性读取） |
| GCM_DecryptFile | `GCM_DecryptFile(inPath, outPath string, key []byte) error` | 解密文件（一次性读取） |

### AES-CTR+HMAC（流式，低内存，带认证）

| 函数 | 签名 | 说明 |
|------|------|------|
| CTRHMAC_Encrypt | `CTRHMAC_Encrypt(plaintext, key []byte) ([]byte, error)` | 加密并认证，输出格式: IV(16B) + HMAC(32B) + ciphertext |
| CTRHMAC_Decrypt | `CTRHMAC_Decrypt(data, key []byte) ([]byte, error)` | 验证 HMAC 后解密，验证失败返回 ErrCorrupted |
| CTRHMAC_EncryptFile | `CTRHMAC_EncryptFile(inPath, outPath string, key []byte) error` | 流式加密并认证文件，低内存占用 |
| CTRHMAC_DecryptFile | `CTRHMAC_DecryptFile(inPath, outPath string, key []byte) error` | 流式验证 HMAC 后解密，低内存占用 |

### 错误变量

| 变量 | 说明 |
|------|------|
| ErrInvalidKey | 密钥长度不合法（必须为 16/24/32 字节） |
| ErrCorrupted | 数据损坏、被篡改或密钥错误 |
| ErrShortData | 输入数据过短 |

## 依赖说明

本包为零依赖，仅使用 Go 标准库。

## 使用示例

```go
import "github.com/lonevle/gokit/crypto"

// 计算字符串 MD5
md5Str := crypto.MD5("hello world")
// 结果: 5eb63bbbe01eeed093cb22bb8f5acdc3

// 计算文件 MD5
fileMD5, err := crypto.FileMD5("/path/to/file.txt")
if err != nil {
	// 处理错误
}

key := []byte("0123456789abcdef0123456789abcdef") // 32 字节密钥

// AES-CTR 流式加密文件（无认证，适合大文件）
err = crypto.CTR_EncryptFile("/path/to/file.txt", "/path/to/file.txt.enc", key)
if err != nil {
	// 处理错误
}

// AES-CTR 流式解密文件
err = crypto.CTR_DecryptFile("/path/to/file.txt.enc", "/path/to/file.txt.dec", key)
if err != nil {
	// 处理错误
}

// AES-GCM 加密文件（自带认证，适合中小文件）
err = crypto.GCM_EncryptFile("/path/to/file.txt", "/path/to/file.txt.gcm", key)
if err != nil {
	// 处理错误
}

// AES-GCM 解密文件
err = crypto.GCM_DecryptFile("/path/to/file.txt.gcm", "/path/to/file.txt.dec", key)
if err != nil {
	// 处理错误
}

// AES-CTR+HMAC 流式加密并认证文件（带认证，适合大文件）
err = crypto.CTRHMAC_EncryptFile("/path/to/file.txt", "/path/to/file.txt.hmac", key)
if err != nil {
	// 处理错误
}

// AES-CTR+HMAC 流式验证并解密文件
err = crypto.CTRHMAC_DecryptFile("/path/to/file.txt.hmac", "/path/to/file.txt.dec", key)
if err != nil {
	// 处理错误
}
```

## 模式选择建议

| 场景 | 推荐模式 | 理由 |
|------|----------|------|
| 大文件加密 | CTR / CTR+HMAC | 流式处理，内存占用低 |
| 中小文件加密 | GCM | 标准推荐，自带认证，API 简洁 |
| 需要完整性认证 | GCM / CTR+HMAC | 自动检测篡改 |
| 仅需要保密性 | CTR | 无认证开销，性能最好 |
