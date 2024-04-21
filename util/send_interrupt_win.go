//go:build windows

package util

import (
	"fmt"
	"syscall"
)

func SendInterrupt(pid int) error {
	dll, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		return fmt.Errorf("LoadDLL: %v", e)
	}
	defer dll.Release()

	p, e := dll.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		return fmt.Errorf("FindProc: %v", e)
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		return fmt.Errorf("GenerateConsoleCtrlEvent: %v", e)
	}
	return nil
}
