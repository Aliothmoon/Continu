//go:build !linux

package build

import (
	"os"
	"os/exec"
)

func Config(c *exec.Cmd) {
	return
}
func Kill(p *os.Process) error {
	return p.Kill()
}
