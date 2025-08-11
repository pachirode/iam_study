package main

import (
	"math/rand"
	"time"

	"github.com/pachirode/iam_study/internal/authzserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	authzserver.NewApp("iam-authz-server").Run()
}
