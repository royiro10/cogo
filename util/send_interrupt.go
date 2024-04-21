//go:build linux || freebsd || darwin

package util

import "syscall"

func SendInterrupt(pid int) error {
	return syscall.Kill(pid, syscall.SIGINT)
}
