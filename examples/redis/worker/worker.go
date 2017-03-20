package main

import (
	_ "github.com/xLegoz/go-swallow/examples"
	"github.com/xLegoz/go-swallow/proto"
	"github.com/xLegoz/go-swallow/workers"
)

func main() {
	// Create basic worker
	w, _ := workers.NewRedisWorker(&proto.RedisWorkerOptions{
		Connection: (&proto.RedisConnectionOptions{}).SetDefaults(),
	})

	// Process messages from a queue
	go w.Process(&proto.RedisWorkerProcessOptions{
		Queue:       "go-swallow-jobs",
		Concurrency: 3,
	})

	// Run worker
	w.Run()
}
