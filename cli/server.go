package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RobertGrantEllis/t9/server"
)

// StartServerFromFlags instantiates and starts a server according to the designated flags. If the flags are invalid
// then the program exits with an error code.
func StartServerFromFlags(args ...string) {
	server := constructServerFromFlags(args...)
	runServerUntilInterrupt(server)
}

func constructServerFromFlags(args ...string) server.Server {
	configuration := constructServerConfigurationFromFlags(args...)

	server, err := server.New(*configuration)
	if err != nil {
		fail(err)
	}

	return server
}

func runServerUntilInterrupt(server server.Server) {
	go stopServerOnInterrupt(server)
	server.Start()
}

func stopServerOnInterrupt(server server.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	stopping := false
	for {
		<-sigChan
		fmt.Println()
		if !stopping {
			go server.Stop()
			stopping = true
		}
	}
}
