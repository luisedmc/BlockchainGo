package main

import (
	"os"
)

func main() {
	defer os.Exit(0)

	cli := CommandLine{}
	cli.Run()
}
