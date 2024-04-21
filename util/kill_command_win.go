//go:build windows

package util

import (
	"os/exec"
	"strconv"
)

func KillCmd(pid int) *exec.Cmd {
	return exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F")
}
