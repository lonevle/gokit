//go:build windows

package iis

import (
	"fmt"
	"strings"

	"github.com/lonevle/gokit/shell"
)

const appcmdPath = `C:\Windows\System32\inetsrv\appcmd.exe`

// runAppCmd 执行 appcmd 命令并返回处理后的输出
func runAppCmd(args ...string) (string, error) {
	if output, err := shell.Exec(appcmdPath, args); err != nil {
		return string(output), err
	} else {
		return string(output), nil
	}
}

// StartPool 启动应用程序池
// 参数:
//   - appPoolName: 应用程序池名称
//
// 返回:
//   - 错误信息
func StartPool(appPoolName string) error {
	output, err := runAppCmd("start", "apppool", fmt.Sprintf("/apppool.name:%s", appPoolName))
	if err != nil {
		return fmt.Errorf("启动应用程序池失败: %w, output: %s", err, output)
	}
	return nil
}

// StopPool 停止应用程序池
// 参数:
//   - appPoolName: 应用程序池名称
//
// 返回:
//   - 错误信息
func StopPool(appPoolName string) error {
	output, err := runAppCmd("stop", "apppool", fmt.Sprintf("/apppool.name:%s", appPoolName))
	if err != nil {
		// 如果已经是停止状态，不返回错误
		if strings.Contains(output, "已停止") {
			return nil
		}
		return fmt.Errorf("停止应用程序池失败: %w, output: %s", err, output)
	}
	return nil
}

// StartSite 启动网站
// 参数:
//   - siteName: 网站名称
//
// 返回:
//   - 错误信息
func StartSite(siteName string) error {
	output, err := runAppCmd("start", "site", siteName)
	if err != nil {
		return fmt.Errorf("启动网站失败: %w, output: %s", err, output)
	}
	return nil
}

// StopSite 停止网站
// 参数:
//   - siteName: 网站名称
//
// 返回:
//   - 错误信息
func StopSite(siteName string) error {
	output, err := runAppCmd("stop", "site", siteName)
	if err != nil {
		return fmt.Errorf("停止网站失败: %w, output: %s", err, output)
	}
	return nil
}
