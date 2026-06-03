# file

文件和路径工具包，提供文件存在检查、编码自动转换、JSON 注释去除等功能。

## 安装

```bash
go get github.com/lonevle/gokit/file
```

## 函数列表

| 函数 | 签名 | 说明 |
|------|------|------|
| Exists | `Exists(path string) (bool, error)` | 检查路径是否存在 |
| ReadFile | `ReadFile(filename string) ([]byte, error)` | 读取文件，自动跳过 BOM 并转换 GBK 到 UTF-8 |
| ReadJson | `ReadJson(filename string) ([]byte, error)` | 读取 JSON 文件，自动处理编码并去除注释 |

## 依赖说明

本包依赖以下外部库：
- `golang.org/x/text`：用于 GBK 编码转换
- `github.com/adhocore/jsonc`：用于去除 JSON 注释

## 使用示例

### 检查文件是否存在

```go
import "github.com/lonevle/gokit/file"

exists, err := file.Exists("/path/to/file.txt")
if err != nil {
    // 处理错误
}
// exists: true/false
```

### 读取文件（自动处理 BOM 和 GBK 编码）

```go
import "github.com/lonevle/gokit/file"

// 读取文件，自动跳过 BOM 头并将 GBK 转换为 UTF-8
content, err := file.ReadFile("data.txt")
if err != nil {
    // 处理错误
}
```

### 读取带注释的 JSON 文件

```go
import "github.com/lonevle/gokit/file"

// 自动处理 BOM、GBK 编码转换并去除 JSON 注释
jsonData, err := file.ReadJson("config.json")
if err != nil {
    // 处理错误
}
```
