package common

import (
	"os"
	"path/filepath"
)

func JoinWithWorkDir(paths ...string) string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exe)
	paths = append([]string{exeDir}, paths...)
	jointPath := filepath.Join(paths...)
	return jointPath
}
