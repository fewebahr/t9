package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RobertGrantEllis/t9/server"
)

func startServerWithArgs(args ...string) {

	server := getServerFromArgs(args...)
	runServerUntilInterrupt(server)
}

func getServerFromArgs(args ...string) server.Server {

	configuration := getServerConfigurationFromArgs(args...)

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

	sigChan := make(chan os.Signal)
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
