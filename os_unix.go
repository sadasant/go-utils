// +build darwin freebsd linux netbsd openbsd

package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// user.Current() requires cgo and thus doesn't work with cross-compiling.
// The following is an alternative that matches how the Heroku Toolbelt
// works, though per @fdr it may not be correct for all cases (when users have
// modified their home dir).
//
// http://stackoverflow.com/questions/7922270/obtain-users-home-directory
func HomePath() string {
	return os.Getenv("HOME")
}

func GetPathFromHome(path string) string {
	parts := strings.Split(path, "/")
	full := HomePath()
	for i, part := range parts {
		full += "/" + part
		if i < len(parts)-1 {
			os.Mkdir(full, os.ModeDir)
		}
	}
	return full
}

func FullPath(path string) string {
	path, _ = filepath.Abs(path)
	return path
}

// CommandPath, attempts to locate a binary
func CommandPath(name string) string {
	path, err := exec.Command("which", name).Output()
	if err != nil {
		return ""
	}
	str_path := string(path)
	str_path = str_path[:len(str_path)-1]

	if FileExist(str_path) {
		return str_path
	}

	return ""
}
