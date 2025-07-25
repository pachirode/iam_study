package log_test

import (
	"fmt"
	"testing"

	"github.com/pachirode/iam_study/pkg/log"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := log.NewOptions()
	opts.AddFlags(pflag.CommandLine)

	args := []string{"--log.level=debug"}
	err := pflag.CommandLine.Parse(args)
	assert.Nil(t, err)
	assert.Equal(t, "debug", opts.Level)
}

func TestValidate(t *testing.T) {
	opts := log.NewOptions()
	opts.Format = "test"
	errs := opts.Validate()
	expected := `[not a valid log format "test"]`
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))
}
