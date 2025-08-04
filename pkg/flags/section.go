package flags

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

type NamedFlagSets struct {
	Order      []string
	FlagSetMap map[string]*pflag.FlagSet
}

func (nfs *NamedFlagSets) GetFlagSet(name string) *pflag.FlagSet {
	if nfs.FlagSetMap == nil {
		nfs.FlagSetMap = map[string]*pflag.FlagSet{}
	}

	if _, ok := nfs.FlagSetMap[name]; !ok {
		nfs.FlagSetMap[name] = pflag.NewFlagSet(name, pflag.ExitOnError)
		nfs.Order = append(nfs.Order, name)
	}

	return nfs.FlagSetMap[name]
}

func PrintSections(w io.Writer, nfs NamedFlagSets, cols int) {
	for _, name := range nfs.Order {
		fs := nfs.FlagSetMap[name]
		if !fs.HasFlags() {
			continue
		}

		wideFS := pflag.NewFlagSet("", pflag.ExitOnError)
		wideFS.AddFlagSet(fs)

		var zzz string
		if cols > 24 {
			zzz = strings.Repeat("z", cols-24)
			wideFS.Int(zzz, 0, zzz)
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "\n%s flags:\n\n%s", strings.ToUpper(name[:1])+name[1:], wideFS.FlagUsagesWrapped(cols))

		if cols > 24 {
			i := strings.Index(buf.String(), zzz)
			lines := strings.Split(buf.String()[:i], "\n")
			fmt.Fprint(w, "%s", strings.Join(lines[:len(lines)-1], "\n"))
			fmt.Fprintln(w)
		} else {
			fmt.Fprint(w, "%s", buf.String())
		}
	}
}
