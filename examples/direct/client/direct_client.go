package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/clients"
	"github.com/xLegoz/go-swallow/jobs"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	client, err := clients.NewDirectClient(&clients.DirectOptions{
		Address: "localhost",
		Port:    8000,
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()
	log.Infof("Connected")

	_, err = client.Perform(
		jobs.Add,
		&jobs.ArithArgs{A: 1, B: 2},
		&clients.EnqueueOptions{Retry: true},
	)
	if err != nil {
		panic(err)
	}
	log.Infof("Sent perform 1")

	_, err = client.Perform(jobs.Add,
		&jobs.ArithArgs{A: 1, B: 2},
		&clients.EnqueueOptions{At: 10}, // 10 minutes
	)
	if err != nil {
		panic(err)
	}
	log.Infof("Sent perform 2")

	resultChannel, err := client.Perform(
		jobs.Add,
		&jobs.ArithArgs{A: 1, B: 2},
	)
	if err != nil {
		panic(err)
	}

	log.Infof("Sent perform 3, blocking")
	result := <-resultChannel
	log.Infof("Got a performed result %v", result)
}
