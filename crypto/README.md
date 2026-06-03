# crypto

加密工具包，提供 MD5 哈希计算功能。

## 安装

```bash
go get github.com/lonevle/gokit/crypto
```

## 函数列表

| 函数 | 签名 | 说明 |
|------|------|------|
| MD5 | `MD5(str string) string` | 计算字符串的 MD5 哈希值，返回 32 位小写字符串 |
| FileMD5 | `FileMD5(path string) (string, error)` | 计算文件的 MD5 哈希值，返回 32 位小写字符串 |

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
```
