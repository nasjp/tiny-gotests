package model

import "fmt"

const (
	ErrStr = "error"
)

type Func struct {
	*Obj
	Receiver   *Obj
	Parameters Objs
	Results    Objs
}

type Obj struct {
	Package   string
	Typ       string
	Name      string
	IsPointer bool
}

type Objs []*Obj

func (o *Obj) Type() string {
	if o == nil {
		return ""
	}

	return o.Typ
}

type Funcs []*Func

func (fs Funcs) Search(funcName string) (*Func, error) {
	for _, f := range fs {
		if f.Name == funcName {
			return f, nil
		}
	}

	return nil, fmt.Errorf("func %s is not found", funcName)
}

func (fs Funcs) String() string {
	var txt string
	for _, f := range fs {
		if txt != "" {
			txt += "\n"
		}

		txt += f.String()
	}

	return txt
}

func (f *Func) String() string {
	return fmt.Sprintf("%s.%s%s%s %s", f.Package, f.receiver(), f.Name, f.parameters(), f.results())
}

func (f *Func) receiver() string {
	if f.Receiver == nil {
		return ""
	}

	return "(" + f.Receiver.str() + ") "
}

func (f *Func) parameters() string {
	var txt string
	for _, p := range f.Parameters {
		if txt != "" {
			txt += ", "
		}

		txt += p.str()
	}

	return "(" + txt + ")"
}

func (f *Func) results() string {
	var (
		txt      string
		cloParen string
	)

	for _, r := range f.Results {
		if txt != "" {
			txt = "(" + txt + ", "
			cloParen = ")"
		}

		txt += r.str() + cloParen
	}

	return txt
}

func (o *Obj) str() string {
	var txt string
	if o.Name != "" {
		txt = o.Name + " "
	}

	if o.Package != "" {
		txt += o.Package + "."
	}

	return txt + o.Typ
}
