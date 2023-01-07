package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/fewebahr/t9/server"
)

func constructServerConfigurationFromFlags(args ...string) *server.Configuration {
	if len(args) == 0 {
		fail(errors.New(`no arguments provided`))
	}

	// creates a flagset and a configuration and sets up the flagset so it parses into the configuration
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	configuration := server.NewConfiguration()

	fs.StringVar(&configuration.LogLevel, `log-level`, configuration.LogLevel, `debug|info|warn|error`)
	fs.StringVar(&configuration.Address, `address`, configuration.Address, `listening address for t9 server`)
	fs.StringVar(&configuration.DictionaryFile, `dictionary`, configuration.DictionaryFile, `dictionary (defaults to built-in English dictionary)`)
	fs.IntVar(&configuration.CacheSize, `cache-size`, configuration.CacheSize, `cache size for t9 words`)
	fs.StringVar(&configuration.CertificateFile, `certificate`, configuration.CertificateFile, `SSL/TLS certificate in PEM format`)
	fs.StringVar(&configuration.KeyFile, `key`, configuration.KeyFile, `SSL/TLS private key in PEM format`)

	addHelpAndParse(fs, args...)

	return &configuration
}

func addHelpAndParse(fs *flag.FlagSet, args ...string) {
	// put in a help flag
	help := false
	fs.BoolVar(&help, `help`, help, `show this help message`)

	// suppress failure output so we can do our own thing
	fs.SetOutput(io.Discard)

	// parse the arguments into the configuration
	err := fs.Parse(args[1:])

	// now we are done parsing, re-enable output writer
	fs.SetOutput(os.Stderr)

	// figure out if we want to print help
	if err != nil && err != flag.ErrHelp {
		printHelp(fs, err)
	} else if fs.NArg() != 0 {
		printHelp(fs, fmt.Errorf(`unexpected argument: '%s'`, fs.Arg(0)))
	} else if help || err == flag.ErrHelp {
		printHelp(fs, nil)
	}
}

func printHelp(fs *flag.FlagSet, err error) {
	exitCode := 0
	if err != nil {
		printError(err)
		fmt.Fprintln(os.Stderr)
		exitCode = 1
	}

	fmt.Fprintf(os.Stderr, "Usage :\n")
	fs.PrintDefaults()

	os.Exit(exitCode)
}
