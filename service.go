package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/takama/daemon"
)

const (
	name        = "revive"
	description = "revive container service test"
	exitMessage = "Service exited"
)

var stdlog, errlog *log.Logger

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	usage := "Usage: cron_job install | remove | start | stop | status "

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	// do healthcheck

	// wait for intterrupt signal
	killSignal := <-interrupt
	stdlog.Println("Got signal: ", killSignal)

	return exitMessage, nil
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}
