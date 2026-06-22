package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// UserHomeDir returns the current user's home directory.
func UserHomeDir() (string, error) {
	return resolveUserHomeDir(runtime.GOOS, os.UserHomeDir, os.Getenv)
}

// ExpandHome expands a leading ~ in path to the current user's home directory.
func ExpandHome(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	home, err := resolveUserHomeDir(runtime.GOOS, os.UserHomeDir, os.Getenv)
	if err != nil {
		return "", err
	}

	return filepath.Join(home, path[1:]), nil
}

func resolveUserHomeDir(goos string, lookupHomeDir func() (string, error), getenv func(string) string) (string, error) {
	if home, err := lookupHomeDir(); err == nil {
		return home, nil
	}

	if goos == "windows" {
		if home := getenv("HOME"); home != "" {
			return home, nil
		}

		drive, path := getenv("HOMEDRIVE"), getenv("HOMEPATH")
		if drive != "" && path != "" {
			return drive + path, nil
		}
	}

	return "", errors.New("home directory is not defined")
}
