package util

import (
	"os"
	"os/signal"
	"syscall"
)

func HandleSignals(quitCallback func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for sig := range signals {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			quitCallback()
		}
	}
}
