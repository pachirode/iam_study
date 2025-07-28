package homedir

import (
	"os"
	"path/filepath"
	"runtime"
)

func HomeDir() string {
	if runtime.GOOS != "windows" {
		return os.Getenv("HOME")
	}
	home := os.Getenv("HOME")
	homeDirveHomePath := ""
	if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
		homeDirveHomePath = homeDrive + homePath
	}
	userProfile := os.Getenv("USERPROFILE")

	for _, p := range []string{home, homeDirveHomePath, userProfile} {
		if len(p) == 0 {
			continue
		}

		if _, err := os.Stat(filepath.Join(p, ".apimachinery", "config")); err != nil {
			continue
		}

		return p
	}

	firstSetPath := ""
	firstExistingPath := ""

	for _, p := range []string{home, userProfile, homeDirveHomePath} {
		if len(p) == 0 {
			continue
		}

		if len(firstSetPath) == 0 {
			firstSetPath = p
		}

		info, err := os.Stat(p)
		if err != nil {
			continue
		}

		if len(firstExistingPath) == 0 {
			firstExistingPath = p
		}

		if info.IsDir() && info.Mode().Perm()&(1<<(uint(7))) != 0 {
			return p
		}
	}

	if len(firstExistingPath) > 0 {
		return firstExistingPath
	}

	if len(firstSetPath) > 0 {
		return firstSetPath
	}

	return ""
}
