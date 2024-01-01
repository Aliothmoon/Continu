//go:build linux

package build

import (
	"os"
	"os/exec"
	"syscall"
)

func Config(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
func Kill(p *os.Process) error {
	return syscall.Kill(-p.Pid, syscall.SYS_KILL)
}
