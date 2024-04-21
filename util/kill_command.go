//go:build linux || freebsd || darwin

package util

import (
	"os/exec"
	"strconv"
)

func KillCmd(pid int) *exec.Cmd {
	return exec.Command("kill", "-9", strconv.Itoa(pid))
}
