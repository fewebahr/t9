package cli

import (
	"strings"

	"github.com/pkg/errors"
)

func StartCLI(args ...string) {

	subcommand := ``
	if len(args) > 1 {
		subcommand = strings.ToLower(args[1])
		args = append([]string{subcommand}, args[2:]...)
	}

	switch subcommand {
	case `server`:
		startServerWithArgs(args...)
	case `client`:
		runClientWithArgs(args...)
	case `-h`:
		fallthrough
	case `-help`:
		fallthrough
	case `--help`:
		fallthrough
	case `help`:
		println(`subcommand is required (server|client)`)
	case ``:
		fail(errors.New(`subcommand is required (server|client)`))
	default:
		fail(errors.Errorf(`unexpected subcommand: '%s'`, subcommand))
	}

}
