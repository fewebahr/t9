package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func printError(err error) {
	fmt.Fprintln(os.Stderr, color.RedString(fmt.Sprintf(`%s: %s`, `error`, err.Error())))
}

func fail(err error) {
	printError(err)
	os.Exit(1)
}
