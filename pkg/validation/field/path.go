package field

import (
	"bytes"
	"fmt"
	"strconv"
)

type Path struct {
	name   string
	index  string
	parent *Path
}

func NewPath(name string, moreName ...string) *Path {
	r := &Path{name: name, parent: nil}
	for _, anotherName := range moreName {
		r = &Path{name: anotherName, parent: r}
	}

	return r
}

func (p *Path) Root() *Path {
	for ; p.parent != nil; p = p.parent {
	}

	return p
}

func (p *Path) Child(name string, moreNames ...string) *Path {
	r := NewPath(name, moreNames...)
	r.Root().parent = p
	return r
}

func (p *Path) Index(index int) *Path {
	return &Path{index: strconv.Itoa(index), parent: p}
}

func (p *Path) Key(key string) *Path {
	return &Path{index: key, parent: p}
}

func (p *Path) String() string {
	elems := []*Path{}
	for ; p != nil; p = p.parent {
		elems = append(elems, p)
	}

	buf := bytes.NewBuffer(nil)
	for i := range elems {
		p := elems[len(elems)-1-i]
		if p.parent != nil && len(p.name) > 0 {
			buf.WriteString(".")
		}
		if len(p.name) > 0 {
			buf.WriteString(p.name)
		} else {
			fmt.Fprintf(buf, "[%s]", p.index)
		}
	}
	return buf.String()
}
