package processor

import (
	"fmt"
	"os"

	"github.com/nasjp/tiny-gotests/internal/generator"
	"github.com/nasjp/tiny-gotests/internal/parser"
)

func Run(filePath string, funcName string) {
	if err := run(filePath, funcName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(filePath string, funcName string) error {
	fs, err := parser.Parse(filePath)
	if err != nil {
		return err
	}

	code := generator.Generate(fs)

	fmt.Println(code)

	return nil
}
