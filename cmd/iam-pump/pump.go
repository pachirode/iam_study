package main

import (
	"math/rand"
	"time"

	"github.com/pachirode/iam_study/internal/pump"
)

func main() {
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	pump.NewApp("iam-pump").Run()
}
