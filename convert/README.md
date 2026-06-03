# convert

编码转换工具包，提供 BOM 头处理和 GBK/UTF-8 编码转换功能。

## 安装

```bash
go get github.com/lonevle/gokit/convert
```

## 函数列表

| 函数 | 签名 | 说明 |
|------|------|------|
| SkipBOM | `SkipBOM(filename string) ([]byte, error)` | 读取文件并跳过 BOM 头 |
| GBKToUTF8 | `GBKToUTF8(data []byte) ([]byte, error)` | 将 GBK 编码转换为 UTF-8 编码 |

## 依赖说明

本包依赖 `golang.org/x/text` 用于 GBK 编码转换。

## 使用示例

### BOM 头处理

```go
import "github.com/lonevle/gokit/convert"

// 读取文件并自动跳过 BOM 头
content, err := convert.SkipBOM("/path/to/file.txt")
if err != nil {
    // 处理错误
}
// content 为剔除 BOM 头后的字节数组
```

支持的 BOM 类型：
- UTF-8: `EF BB BF`
- UTF-16 大端序: `FE FF`
- UTF-16 小端序: `FF FE`
- UTF-32 大端序: `00 00 FE FF`
- UTF-32 小端序: `FF FE 00 00`

### GBK 转 UTF-8

```go
import "github.com/lonevle/gokit/convert"

// 将 GBK 编码转换为 UTF-8
utf8Data, err := convert.GBKToUTF8(gbkData)
if err != nil {
    // 处理错误
}

// 如果输入已经是 UTF-8 编码，会直接返回原数据
```
