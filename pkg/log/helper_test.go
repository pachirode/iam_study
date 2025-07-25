package log_test

import (
	"testing"

	"github.com/pachirode/iam_study/pkg/log"
)

func TestWithName(t *testing.T) {
	defer log.Flush()

	logger := log.WithName("test")
	logger.Infow("Infow name", "foo", "bar")
}

func TestWithValues(t *testing.T) {
	defer log.Flush()

	logger := log.WithValues("key", "value")
	logger.Info("Test values")
}

func TestV(t *testing.T) {
	defer log.Flush()

	log.V(0).Infow("Test V ", "foo", "bar")
	log.V(1).Infow("Test V ", "foo", "bar")
}
