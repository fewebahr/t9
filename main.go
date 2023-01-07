package main

import (
	"os"

	"github.com/fewebahr/t9/cli"
)

func main() {
	cli.StartServerFromFlags(os.Args...)
}
