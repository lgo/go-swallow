package main

import (
	_ "github.com/xLegoz/go-swallow/jobs"
	"github.com/xLegoz/go-swallow/master/client"
	"github.com/xLegoz/go-swallow/master/worker"
)

func main() {

	// TODO: Start HTTP server with dashboards
	go client.Listen(8000)
	worker.Listen(8080)
}
