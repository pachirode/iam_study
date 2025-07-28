package app

import (
	"os"
	"runtime"
	"strings"

	"github.com/pachirode/iam_study/pkg/log"
)

func FormatBasename(basename string) string {
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}
