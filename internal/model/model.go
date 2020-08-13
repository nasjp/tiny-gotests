package model

import "fmt"

type Func struct {
	Package    string
	Name       string
	IsExported bool
	Receiver   string
}

type Funcs []*Func

func (fs Funcs) String() string {
	var txt string
	for _, f := range fs {
		if txt != "" {
			txt += "\n"
		}

		txt += fmt.Sprintf("package: %s, name: %s, isExported: %t", f.Package, f.Name, f.IsExported)
		if f.Receiver != "" {
			txt += fmt.Sprintf(", receiver: %s", f.Receiver)
		}
	}

	return txt
}
