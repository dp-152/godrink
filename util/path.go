package util

import (
	"path/filepath"
	"runtime"
)

var rootDir string

func init() {
	if rootDir == "" {
		_, currentFile, _, _ := runtime.Caller(0)
		rootDir = filepath.Join(filepath.Dir(currentFile), "..")
	}
}

func RootPath() string {
	return rootDir
}

func OwnPath() string {
	_, callerFile, _, _ := runtime.Caller(1)
	return filepath.Dir(callerFile)
}

func OwnRelPath() string {
	_, callerFile, _, _ := runtime.Caller(1)
	path, err := filepath.Rel(rootDir, filepath.Dir(callerFile))

	if err != nil {
		return ""
	}
	return path
}

func JoinFromRoot(segments ...string) string {
	return filepath.Join(append([]string{rootDir}, segments...)...)
}
