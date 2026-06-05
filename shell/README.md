# shell

Windows 命令行调用工具包，提供对 PowerShell 和 CMD 的封装。

## 安装

```bash
go get github.com/lonevle/gokit/shell
```

## 功能说明

- 执行 PowerShell 命令或脚本文件
- 执行 CMD 命令或批处理脚本文件
- 自动将 GBK 输出转换为 UTF-8
- 支持自定义错误检测（补充错误文本匹配）

## 函数列表

### shell
```go
shell.Exec("xxx.exe", []string{"arg"})
```


### PowerShell

| 函数 | 签名 | 说明 |
|------|------|------|
| RunPowerShellCommand | `func RunPowerShellCommand(psCmd string, supplemental []string, config *PSConfig) ([]byte, error)` | 执行 PowerShell 命令字符串 |
| RunPowerShellFile | `func RunPowerShellFile(path string, args []string, supplemental []string, config *PSConfig) ([]byte, error)` | 执行 PowerShell 脚本文件 |

### CMD

| 函数 | 签名 | 说明 |
|------|------|------|
| RunCmdCommand | `func RunCmdCommand(cmd string, config *CmdConfig) ([]byte, error)` | 执行 CMD 命令字符串 |
| RunCmdScript | `func RunCmdScript(path string, args []string, config *CmdConfig) ([]byte, error)` | 执行 CMD 批处理脚本文件 |

## 依赖说明

- `github.com/lonevle/gokit/convert`：GBK/UTF-8 编码转换

## 使用示例

### PowerShell 示例

```go
package main

import (
    "fmt"
    "log"

    "github.com/lonevle/gokit/shell"
)

func main() {
    // 执行简单命令
    out, err := shell.RunPowerShellCommand("Get-Date", nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))

    // 使用 PowerShell 7 并检测特定错误文本
    config := &shell.PSConfig{
        ErrAction: shell.Stop,
        UsePwsh7:  true,
    }
    out, err = shell.RunPowerShellCommand("Get-Process", []string{"not found"}, config)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))

    // 执行脚本文件
    out, err = shell.RunPowerShellFile(`C:\scripts\test.ps1`, []string{"-Verbose"}, nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))
}
```

### CMD 示例

```go
package main

import (
    "fmt"
    "log"

    "github.com/lonevle/gokit/shell"
)

func main() {
    // 执行简单命令
    out, err := shell.RunCmdCommand("dir /b", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))

    // 执行批处理脚本
    out, err = shell.RunCmdScript(`C:\scripts\test.bat`, []string{"arg1", "arg2"}, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(out))
}
```
