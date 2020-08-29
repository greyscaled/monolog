// monolog is a simple CLI to simplify reading daily git logs from a monorepo.
// The program uses a configuration file located at ~/.monologconfig with read
// and write permissions.
package main

import (
	"fmt"
	"os"
	"strings"
	"vapurrmaid/monolog/xgit"
)

func main() {
	l := len(os.Args)

	switch l {
	case 2:
		xgit.LogLatest(os.Args[1])
	case 1:
		xgit.LogLatest(".")
	default:
		usage()
		os.Exit(0)
	}

	os.Exit(0)
}

func usage() {
	args := []string{
		os.Args[0],
		"[path]",
	}
	fmt.Printf("Usage:\t%s\n", strings.Join(args, " "))
	fmt.Printf("\t%s\n", "path defaults to the current directory")
}
