package clients

import (
	"bufio"
	"bytes"

	log "github.com/Sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/xLegoz/go-swallow/encoders"
	"github.com/xLegoz/go-swallow/proto"
	"github.com/xLegoz/go-swallow/util"
)

type RedisClient struct {
	queue string
	pool  *redis.Pool
}

func NewRedisClient(options *proto.RedisClientOptions) (redisClient *RedisClient, err error) {
	log.WithFields(log.Fields{
		"address": options.Connection.Address,
		"queue":   options.Queue,
	}).Info("Initialized Redis client")
	redisClient = &RedisClient{
		queue: options.Queue,
		pool:  util.GetRedisPool(options.Connection),
	}
	return
}

func (c *RedisClient) Perform(fn interface{}, args interface{}, extra ...interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()

	// Create job from function and arguments
	jobRequest, err := util.CreateJobRequest(fn, args)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error creating job request")
		return err
	}

	// Generate and serialize message
	msg := &proto.Message{
		JobID:  util.GenerateJobID(),
		Job:    jobRequest,
		Status: 1,
	}
	var buff bytes.Buffer
	writer := bufio.NewWriter(&buff)
	err = encoders.GobEncode(writer, msg)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error encoding message")
		return err
	}
	writer.Flush()
	data := buff.Bytes()

	// Push job to Redis
	_, err = conn.Do("LPUSH", c.queue, data)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  data,
		}).Error("Error queueing message")
		return err
	}
	log.WithFields(log.Fields{
		"jobID":               msg.JobID,
		"jobRequest.FuncName": jobRequest.FuncName,
		"jobRequest.FuncArgs": jobRequest.FuncArgs,
	}).Info("Queued job")

	return nil
}

func (c *RedisClient) Close() error {
	return c.pool.Close()
}
