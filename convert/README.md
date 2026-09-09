# convert

编码识别与转换工具包，将任意编码（UTF-8 / GBK / UTF-16 / UTF-32）统一转为 UTF-8。只做纯转换，不进行文件 IO。

## 安装

```bash
go get github.com/lonevle/gokit/convert
```

## 函数列表

| 函数 | 签名 | 说明 |
|------|------|------|
| ToUTF8 | `ToUTF8(data []byte) ([]byte, error)` | 识别编码并统一转为 UTF-8 |
| NewUTF8Reader | `NewUTF8Reader(r io.Reader) (io.Reader, error)` | 流式版本，内存 O(缓冲区)，适合大文件 |

## 依赖说明

本包依赖 `golang.org/x/text`。

## 使用示例

### 内存转换

```go
import "github.com/lonevle/gokit/convert"

// GBK / UTF-16 / 带 BOM 等编码自动识别并转为 UTF-8
u8, err := convert.ToUTF8(data)
if err != nil {
    // 无法识别的编码（如二进制）
}
```

支持的编码：

- UTF-8（含 BOM，自动剥离）
- GBK（Windows 控制台默认 cp936）
- UTF-16 大端序: `FE FF`，小端序: `FF FE`
- UTF-32 大端序: `00 00 FE FF`，小端序: `FF FE 00 00`

### 大文件流式转换

```go
import (
    "io"
    "os"

    "github.com/lonevle/gokit/convert"
)

f, err := os.Open("/path/to/big.log")
if err != nil {
    // 处理错误
}
defer f.Close()

r, err := convert.NewUTF8Reader(f)
if err != nil {
    // 处理错误
}
io.Copy(os.Stdout, r)
```

## 已知局限

无 BOM 的 GBK 内容若恰好构成合法 UTF-8，会被误判为 UTF-8。
