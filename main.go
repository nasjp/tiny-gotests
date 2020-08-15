package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nasjp/tiny-gotests/internal/processor"
)

func main() {
	var (
		filePath = flag.String("p", "", `target file path`)
		funcName = flag.String("f", "", `target func name`)
	)

	flag.Parse()
	run(*filePath, *funcName)
}

func run(filePath string, funcName string) {
	code, err := processor.Run(filePath, funcName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(code))

	os.Exit(0)
}
