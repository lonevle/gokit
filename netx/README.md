# netx

网络工具包，提供 IP 地址相关的工具函数。

## 安装

```bash
go get github.com/lonevle/gokit/netx
```

## 函数列表

| 函数 | 签名 | 说明 |
|------|------|------|
| IsPrivateIP | `IsPrivateIP(addr string) bool` | 判断 IP 是否为私有/内网地址，包含 RFC 1918 私有地址、回环地址和链路本地地址 |

## 支持的地址类型

- **RFC 1918 私有地址**: 10.0.0.0/8、172.16.0.0/12、192.168.0.0/16、fc00::/7
- **回环地址**: 127.0.0.0/8、::1
- **链路本地地址**: 169.254.0.0/16、fe80::/10

## 依赖说明

本包为零依赖，仅使用 Go 标准库。

## 使用示例

```go
import "github.com/lonevle/gokit/netx"

// 判断是否为私有/内网 IP
netx.IsPrivateIP("192.168.1.1")  // true
netx.IsPrivateIP("10.0.0.1")     // true
netx.IsPrivateIP("127.0.0.1")    // true (回环地址)
netx.IsPrivateIP("169.254.1.1")  // true (链路本地地址)
netx.IsPrivateIP("8.8.8.8")      // false (公网地址)
```
