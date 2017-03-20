package main

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/xLegoz/go-swallow/jobs"
	"github.com/xLegoz/go-swallow/proto"
	"github.com/xLegoz/go-swallow/workers"
)

type myMiddleware struct{}

func (r *myMiddleware) Call(queue string, message *proto.Message, next func() bool) (acknowledge bool) {
	log.Debug("Do something before a message is processed")
	acknowledge = next()
	log.Debug("Do something after a message is processed")
	return
}

func main() {
	w, err := workers.NewDirectWorker(&proto.DirectWorkerOptions{
		Port: 8000,
	})
	if err != nil {
		panic(err)
	}

	w.Middleware.Append(&myMiddleware{})

	// Process everything
	go w.Process()

	w.Run()
}
