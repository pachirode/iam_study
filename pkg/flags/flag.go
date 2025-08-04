package flags

import (
	goflag "flag"
	"strings"

	"github.com/spf13/pflag"

	"github.com/pachirode/iam_study/pkg/log"
)

func SepNormalizeNameFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		newName := strings.ReplaceAll(name, "_", "-")
		log.Warnf("%s is be replaced by %s", name, newName)

		return pflag.NormalizedName(newName)
	}

	return pflag.NormalizedName(name)
}

func InitFlags() {
	pflag.CommandLine.SetNormalizeFunc(SepNormalizeNameFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		log.Debugf("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
