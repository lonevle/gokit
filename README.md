# gokit

Go 工具库，提供常用的加密、字符串、文件等工具函数。

## 安装

```bash
go get github.com/lonevle/gokit
```

## 功能模块

| 模块 | 路径 | 说明 | 依赖 |
|------|------|------|------|
| [crypto](crypto/README.md) | `github.com/lonevle/gokit/crypto` | MD5 哈希计算、AES 文件加密/解密 | 零依赖 |
| [file](file/README.md) | `github.com/lonevle/gokit/file` | 文件和路径工具 | `golang.org/x/text`、`github.com/adhocore/jsonc` |
| [convert](convert/README.md) | `github.com/lonevle/gokit/convert` | 编码转换（BOM 处理、GBK/UTF-8 转换） | `golang.org/x/text` |
| [iis](iis/README.md) | `github.com/lonevle/gokit/iis` | IIS 管理（应用池/网站启停） | 零依赖 |
| [netx](netx/README.md) | `github.com/lonevle/gokit/netx` | 网络工具（IP 地址判断） | 零依赖 |