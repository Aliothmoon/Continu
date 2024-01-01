//go:build linux

package build

import (
	"github.com/Aliothmoon/Continu/internal/logger"
	"os"
	"os/exec"
	"syscall"
)

func init() {
	logger.Debug("Build By Linux Kill")
}
func Config(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
func Kill(p *os.Process) error {
	return syscall.Kill(-p.Pid, syscall.SYS_KILL)
}
