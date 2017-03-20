package main

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/clients"
	"github.com/xLegoz/go-swallow/examples"
	"github.com/xLegoz/go-swallow/proto"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 2 {
		log.Fatal("This program requires 2 args")
	}

	first, err := strconv.Atoi(argsWithoutProg[0])
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Arguments need to be valid integers")
	}
	second, err := strconv.Atoi(argsWithoutProg[1])
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Arguments need to be valid integers")
	}

	client, err := clients.NewRedisClient(&proto.RedisClientOptions{
		Connection: proto.DefaultRedisConnectionOptions,
		Queue:      "go-swallow-jobs",
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create Redis client")
	}
	defer client.Close()

	err = client.Perform(
		examples.Add,
		&examples.ArithArgs{A: first, B: second},
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to perform job")
	}

	log.Info("Client finished")
}
