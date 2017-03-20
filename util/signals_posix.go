// +build !windows

package util

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func HandleSignals(quitCallback func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)

	for sig := range signals {
		log.WithFields(log.Fields{
			"signal": sig,
		}).Debug("Received signal")
		switch sig {
		case syscall.SIGINT, syscall.SIGUSR1, syscall.SIGTERM:
			quitCallback()
		}
	}
}
