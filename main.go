package main

import (
	"os"

	"github.com/RobertGrantEllis/t9/cli"
)

func main() {
	cli.StartServerFromFlags(os.Args...)
}
