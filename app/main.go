package main

import (
	"fmt"
	"os"
)

// TODO: Signal Interrupt not working
// TODO: Connections time out after leaving server
// TODO: React to local network scanning

const DEBUG = false

func main() {

	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Fprint(os.Stderr, "Please provide config path as argument\n")
		os.Exit(1)
	}

	srv, err := ServerFromConfig(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config from \"%s\": %s\n", os.Args[1], err.Error())
		os.Exit(2)
	}

	// React to process signals
	// signals := make(chan os.Signal, 1)
	// signal.Notify(signals, os.Interrupt, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	sig := <-signals
	// 	fmt.Fprintf(os.Stdout, "Signal received: %s\n", sig.String())
	// 	srv.Active = false
	// }()

	srv.Listen()
}
