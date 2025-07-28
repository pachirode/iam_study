package flags

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

func AddGlobalFlags(flagSet *pflag.FlagSet, name string) {
	flagSet.BoolP("help", "h", false, fmt.Sprintf("help for %s", name))
}

func RegisterGlobal(local *pflag.FlagSet, globalName string) {
	if flagSet := flag.CommandLine.Lookup(globalName); flagSet != nil {
		pflagGFlag := pflag.PFlagFromGoFlag(flagSet)
		pflagGFlag.Name = strings.ReplaceAll(pflagGFlag.Name, "_", "-")
		local.AddFlag(pflagGFlag)
	} else {
		panic(fmt.Sprintf("failed to find flag in global flagSet (flag): %s", globalName))
	}
}
