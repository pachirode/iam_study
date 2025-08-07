package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/pachirode/iam_study/internal/apiserver"
)

func main() {
	rand.NewSource(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	apiserver.NewApp("iam-apiserver").Run()
}
