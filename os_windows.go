// +build windows

package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const netrcFilename = "_netrc"

func HomePath() string {
	return os.Getenv("HOME")
}

func GetPathFromHome(path string) string {
	// Note: Apparently there' s no MkdirAll in the windows version of Go
	parts := strings.Split(path, "/")
	full := HomePath()
	for i, part := range parts {
		full += "\\" + part
		if i < len(parts)-1 {
			os.Mkdir(full, os.ModeDir)
		}
	}
	return full
}

func FullPath(path string) string {
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "C:\\", "C:\\cygwin\\", 1)
	return path
}

// CommandPath, attempts to locate a binary
func CommandPath(name string) string {
	path, err := exec.Command("which", name).Output()
	if err != nil {
		return ""
	}

	// We'll get a cygwin path, let's make it a windows path
	str_path := string(path)
	str_path = str_path[:len(str_path)-1]
	str_path = strings.Replace(str_path, "\\", "", -1)
	str_path = strings.Replace(str_path, "/cygdrive/c/", "C:\\", -1)
	str_path = strings.Replace(str_path, "/", "\\", -1)

	if _, err := os.Stat(str_path); err == nil {
		return str_path
	}

	return ""
}
