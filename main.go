package main

import (
	"flag"

	"github.com/nasjp/tiny-gotests/internal/processor"
)

func main() {
	var (
		filePath = flag.String("p", "", `target file path`)
		funcName = flag.String("f", "", `target func name`)
	)

	flag.Parse()
	processor.Run(*filePath, *funcName)
}
