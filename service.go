package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
	"github.com/takama/daemon"
)

const (
	name        = "res9"
	description = "revive container service"
	exitMessage = "Service exited"
)

// Service is daemon service
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: res9 install | remove | start | stop | status "

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
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

	// ---------- healthcheck --------------
	c := cron.New()
	// execute rescue func every minute
	c.AddFunc("0 * * * * *", rescue)
	c.Start()
	// wait for intterrupt signal
	killSignal := <-interrupt
	stdlog.Println("Got signal: ", killSignal)

	return exitMessage, nil
}
