package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	listen = flag.String("listen", ":8354", "server listen address")
	runner = flag.Bool("runner", false, "run as a runner")
)

func startRunner() {
	fmt.Println("Starting runner...")
	runFile, err := os.Executable()
	if err != nil {
		log.Printf("E! Failed to get executable file. %s\n", err.Error())
		os.Exit(1)
	}

	for {
		cmd := exec.Command(runFile, "-listen", *listen)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		log.Printf("I! Restarting runner...\n")
		time.Sleep(time.Second * 5)
	}
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if *runner {
		startRunner()
	} else {
		svr := Run()
		interrupt := make(chan os.Signal, 1)
		shutdownSignals := []os.Signal{
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
		}
		signal.Notify(interrupt, shutdownSignals...)

		select {
		case <-interrupt:
			log.Println("Shutting down the server...")
			svr.GracefulStop()
		case <-ctx.Done():
		}

	}

}
