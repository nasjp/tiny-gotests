package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/nasjp/tiny-gotests/internal/model"
)

func Parse(filePath string) (model.Funcs, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)

	if err != nil {
		return nil, fmt.Errorf("target parser.ParseFile(): %v", err)
	}

	var fs model.Funcs

	for _, d := range f.Decls {
		fDecl, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fs = append(fs, &model.Func{
			Package:    f.Name.String(),
			Name:       fDecl.Name.String(),
			IsExported: fDecl.Name.IsExported(),
			Receiver:   pluckReceiver(fDecl),
		})
	}

	return fs, nil
}

func pluckReceiver(fDecl *ast.FuncDecl) string {
	if fDecl.Recv == nil {
		return ""
	}

	return fmt.Sprint(fDecl.Recv.List[0].Type)
}
