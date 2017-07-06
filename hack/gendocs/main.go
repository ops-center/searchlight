package main

import (
	"fmt"
	"log"
	"os"

	"github.com/appscode/go/runtime"
	"github.com/appscode/searchlight/plugins/hyperalert"
	"github.com/spf13/cobra/doc"
)

// ref: https://github.com/spf13/cobra/blob/master/doc/md_docs.md
func main() {
	rootCmd := hyperalert.NewCmd()
	dir := runtime.GOPath() + "/src/github.com/appscode/searchlight/docs/hyperalert"
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	doc.GenMarkdownTree(rootCmd, dir)
}
