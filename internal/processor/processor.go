package processor

import (
	"github.com/nasjp/tiny-gotests/internal/generator"
	"github.com/nasjp/tiny-gotests/internal/parser"
)

func Run(filePath string, funcName string) ([]byte, error) {
	fs, err := parser.Parse(filePath)
	if err != nil {
		return nil, err
	}

	f, err := fs.Search(funcName)
	if err != nil {
		return nil, err
	}

	code, err := generator.Generate(filePath, f)
	if err != nil {
		return nil, err
	}

	return code, nil
}
