# iis

IIS 管理工具包，提供 Windows IIS 应用程序池和网站的启停控制功能。

**注意**：本包仅限 Windows 平台使用，使用了 `//go:build windows` 构建标签。

## 安装

```bash
go get github.com/lonevle/gokit/iis
```

## 函数列表

| 函数 | 签名 | 说明 |
|------|------|------|
| StartPool | `StartPool(appPoolName string) error` | 启动指定名称的应用程序池 |
| StopPool | `StopPool(appPoolName string) error` | 停止指定名称的应用程序池 |
| StartSite | `StartSite(siteName string) error` | 启动指定名称的网站 |
| StopSite | `StopSite(siteName string) error` | 停止指定名称的网站 |

## 依赖说明

本包为零依赖，仅使用 Go 标准库。

## 使用示例

```go
import "github.com/lonevle/gokit/iis"

// 启动应用程序池
err := iis.StartPool("MyAppPool")

// 停止应用程序池（如果已经是停止状态不会报错）
err := iis.StopPool("MyAppPool")

// 启动网站
err := iis.StartSite("MySite")

// 停止网站
err := iis.StopSite("MySite")
```
