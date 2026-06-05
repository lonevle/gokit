//go:build windows
// +build windows

package shell

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	powerShellExe  = `powershell.exe`
	powerShell7Exe = `pwsh.exe`
)

// ErrorAction 定义 PowerShell 的错误处理策略
type ErrorAction string

const (
	// SilentlyContinue 遇到非终止错误时静默继续
	SilentlyContinue ErrorAction = "SilentlyContinue"
	// Continue 遇到错误时继续执行，但会在输出中显示错误
	Continue ErrorAction = "Continue"
	// Stop 将所有 cmdlet 错误视为终止错误
	Stop ErrorAction = "Stop"
	// Inquire 显示调试信息并询问用户是否继续
	Inquire ErrorAction = "Inquire"
)

// PSConfig 表示 PowerShell 调用的配置参数
type PSConfig struct {
	ErrAction ErrorAction // 错误处理策略
	Params    []string    // 额外的 powershell.exe 参数
	UsePwsh7  bool        // 是否使用 PowerShell 7 (pwsh.exe)
}

var (
	// defaultConfig 是 PowerShell 调用的默认配置
	// 默认使用 -NoProfile（不加载用户配置文件）和 -NonInteractive（非交互模式）
	// 并将错误策略设为 Stop
	defaultConfig = PSConfig{
		ErrAction: Stop,
		Params:    []string{"-NoProfile", "-NonInteractive"},
		UsePwsh7:  false,
	}

	// ErrPowerShell 表示 PowerShell 执行返回的错误
	ErrPowerShell = errors.New("powershell error")
	// ErrSupplemental 表示输出中包含用户指定的补充错误文本
	ErrSupplemental = errors.New("supplemental error")
	// ErrUnsupported 表示调用了不支持的功能
	ErrUnsupported = errors.New("unsupported function call")

	errCompile = errors.New("compile error")
)

// supplementalErr 检查 PowerShell 输出中是否包含用户指定的补充错误文本
func supplementalErr(output []byte, supplemental []string) error {
	for _, supplement := range supplemental {
		expression := fmt.Sprintf(`(?i)%s[\s\S]+`, regexp.QuoteMeta(supplement))
		regex, err := regexp.Compile(expression)
		if err != nil {
			return fmt.Errorf("regexp.Compile(%q) returned %v: %w", expression, err, errCompile)
		}
		if regex.Match(output) {
			return fmt.Errorf("output %s contains supplemental error text %q: %w", output, supplement, ErrPowerShell)
		}
	}
	return nil
}

// executePowerShell 执行 PowerShell 并检查补充错误
func executePowerShell(params []string, supplemental []string, pwsh7 bool) ([]byte, error) {
	exe := powerShellExe
	if pwsh7 {
		exe = powerShell7Exe
	}
	out, err := Exec(exe, params)
	if err != nil {
		return out, fmt.Errorf("powershell returned %v: %w", err, ErrPowerShell)
	}
	if err := supplementalErr(out, supplemental); err != nil {
		return out, fmt.Errorf("supplementalErr returned %v: %w", err, ErrSupplemental)
	}
	return out, nil
}

// normalizePSConfig 如果 config 为 nil 则返回默认配置的副本
func normalizePSConfig(config *PSConfig) *PSConfig {
	if config == nil {
		c := defaultConfig
		return &c
	}
	return config
}

// RunPowerShellCommand 执行 PowerShell 命令字符串
// psCmd: 要执行的 PowerShell 命令
// supplemental: 用于检测输出中是否包含特定错误文本的关键词列表
// config: PowerShell 配置，传 nil 使用默认配置
func RunPowerShellCommand(psCmd string, supplemental []string, config *PSConfig) ([]byte, error) {
	cfg := normalizePSConfig(config)
	cmd := fmt.Sprintf(`$ErrorActionPreference="%s"; %s`, cfg.ErrAction, psCmd)
	params := append(cfg.Params, "-Command", cmd)
	return executePowerShell(params, supplemental, cfg.UsePwsh7)
}

// RunPowerShellFile 执行 PowerShell 脚本文件
// path: 脚本文件路径
// args: 传递给脚本的参数
// supplemental: 用于检测输出中是否包含特定错误文本的关键词列表
// config: PowerShell 配置，传 nil 使用默认配置
func RunPowerShellFile(path string, args []string, supplemental []string, config *PSConfig) ([]byte, error) {
	cfg := normalizePSConfig(config)
	params := append(cfg.Params, "-File", path)
	params = append(params, args...)
	return executePowerShell(params, supplemental, cfg.UsePwsh7)
}
