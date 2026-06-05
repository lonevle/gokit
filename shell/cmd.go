//go:build windows
// +build windows

package shell

import (
	"fmt"
	"strings"
)

// CmdConfig 表示 CMD 调用的配置参数
type CmdConfig struct {
	Params []string // 额外的 cmd.exe 参数（如 /C、/K 等）
}

var (
	// defaultCmdConfig 是 CMD 调用的默认配置，默认使用 /C 执行命令后关闭窗口
	defaultCmdConfig = CmdConfig{
		Params: []string{"/C"},
	}
)

// normalizeCmdConfig 如果 config 为 nil 则返回默认配置的副本
func normalizeCmdConfig(config *CmdConfig) *CmdConfig {
	if config == nil {
		c := defaultCmdConfig
		return &c
	}
	return config
}

// runCmd 执行 CMD 命令并返回输出
func runCmd(params []string) ([]byte, error) {
	return Exec("cmd.exe", params)
}

// RunCmdCommand 执行 CMD 命令字符串
// cmd: 要执行的命令（如 "dir /b"）
// config: CMD 配置，传 nil 使用默认配置（/C）
func RunCmdCommand(cmd string, config *CmdConfig) ([]byte, error) {
	cfg := normalizeCmdConfig(config)
	params := append(cfg.Params, cmd)
	return runCmd(params)
}

// RunCmdScript 执行 CMD 批处理脚本文件
// path: 脚本文件路径（如 .bat 或 .cmd 文件）
// args: 传递给脚本的参数
// config: CMD 配置，传 nil 使用默认配置（/C）
func RunCmdScript(path string, args []string, config *CmdConfig) ([]byte, error) {
	cfg := normalizeCmdConfig(config)
	// 使用 call 命令调用脚本，确保脚本执行完毕后返回正确环境
	var cmd string
	if len(args) > 0 {
		cmd = fmt.Sprintf("call %s %s", path, strings.Join(args, " "))
	} else {
		cmd = fmt.Sprintf("call %s", path)
	}
	params := append(cfg.Params, cmd)
	return runCmd(params)
}
