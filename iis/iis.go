//go:build windows

package iis

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/lonevle/gokit/convert"
)

const appcmdPath = `C:\Windows\System32\inetsrv\appcmd.exe`

// runAppCmd 执行 appcmd 命令并返回处理后的输出
func runAppCmd(args ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(appcmdPath, args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	utf8Bytes, _ := convert.GBKToUTF8(outBuf.Bytes())
	stdout = string(utf8Bytes)
	stderr = errBuf.String()
	return
}

// StartPool 启动应用程序池
// 参数:
//   - appPoolName: 应用程序池名称
//
// 返回:
//   - 错误信息
func StartPool(appPoolName string) error {
	_, stderr, err := runAppCmd("start", "apppool", fmt.Sprintf("/apppool.name:%s", appPoolName))
	if err != nil {
		return fmt.Errorf("启动应用程序池失败: %w, stderr: %s", err, stderr)
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
	stdout, stderr, err := runAppCmd("stop", "apppool", fmt.Sprintf("/apppool.name:%s", appPoolName))
	if err != nil {
		// 如果已经是停止状态，不返回错误
		if strings.Contains(stdout, "已停止") {
			return nil
		}
		return fmt.Errorf("停止应用程序池失败: %w, stderr: %s", err, stderr)
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
	_, stderr, err := runAppCmd("start", "site", siteName)
	if err != nil {
		return fmt.Errorf("启动网站失败: %w, stderr: %s", err, stderr)
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
	_, stderr, err := runAppCmd("stop", "site", siteName)
	if err != nil {
		return fmt.Errorf("停止网站失败: %w, stderr: %s", err, stderr)
	}
	return nil
}
