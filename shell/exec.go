package shell

import (
	"fmt"
	"os/exec"

	"github.com/lonevle/gokit/convert"
)

// Exec 执行指定的命令并返回输出, 自动识别编码并统一转为 UTF-8（GBK / UTF-16 等）
// name: 可执行文件名称或路径（如 "cmd.exe" 或 "powershell.exe"）
// params: 传递给可执行文件的参数列表
func Exec(name string, params []string) ([]byte, error) {
	out, err := exec.Command(name, params...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(`exec.Command(%q, %v) command returned: %q: %w`, name, params, out, err)
	}
	return convert.ToUTF8(out)
}
