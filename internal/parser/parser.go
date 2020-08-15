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
			Obj: &model.Obj{
				Package: f.Name.String(),
				Name:    fDecl.Name.String(),
			},
			Receiver:   pluckReceiver(fDecl, f.Name.String()),
			Parameters: pluckParameters(fDecl, f.Name.String()),
			Results:    pluckResults(fDecl, f.Name.String()),
		})
	}

	return fs, nil
}

func pluckReceiver(fDecl *ast.FuncDecl, pkgName string) *model.Obj {
	if fDecl.Recv == nil {
		return nil
	}

	return pluckObject(fDecl.Recv.List[0], pkgName)
}

func pluckParameters(fDecl *ast.FuncDecl, pkgName string) model.Objs {
	if fDecl.Type.Params == nil {
		return nil
	}

	params := make(model.Objs, 0, len(fDecl.Type.Params.List))

	for _, p := range fDecl.Type.Params.List {
		params = append(params, pluckObject(p, pkgName))
	}

	return params
}

func pluckResults(fDecl *ast.FuncDecl, pkgName string) model.Objs {
	if fDecl.Type.Results == nil {
		return nil
	}

	results := make(model.Objs, 0, len(fDecl.Type.Results.List))

	for _, r := range fDecl.Type.Results.List {
		results = append(results, pluckObject(r, pkgName))
	}

	return results
}

func pluckObject(param *ast.Field, pkgName string) *model.Obj {
	var name string
	if len(param.Names) != 0 {
		name = param.Names[0].Name
	}

	switch typ := param.Type.(type) {
	case *ast.Ident:
		// basic type or internal defined value type
		if typ.Obj == nil {
			return &model.Obj{
				Package: "",
				Typ:     typ.Name,
				Name:    name,
			}
		}

		return &model.Obj{
			Package: pkgName,
			Typ:     typ.Name,
			Name:    name,
		}

	case *ast.SelectorExpr:
		internalType := typ.X.(*ast.Ident)

		return &model.Obj{
			Package: internalType.Name,
			Typ:     typ.Sel.Name,
			Name:    name,
		}
	case *ast.StarExpr:
		internalType, ok := typ.X.(*ast.Ident)

		if ok {
			// internal defined pointer type
			return &model.Obj{
				Package:   pkgName,
				Typ:       internalType.Name,
				Name:      name,
				IsPointer: true,
			}
		}

		externalType := typ.X.(*ast.SelectorExpr)

		// external defined pointer type
		return &model.Obj{
			Package:   externalType.X.(*ast.Ident).Name,
			Typ:       externalType.Sel.Name,
			Name:      name,
			IsPointer: true,
		}
	}

	return nil
}
