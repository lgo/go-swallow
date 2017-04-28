package workers

import (
	"bytes"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/xLegoz/go-swallow/encoders"
	"github.com/xLegoz/go-swallow/proto"
	"github.com/xLegoz/go-swallow/util"
)

type RedisWorker struct {
	pool *redis.Pool
	*Worker
}

func NewRedisWorker(options *proto.RedisWorkerOptions) (redisWorker *RedisWorker, err error) {
	redisWorker = &RedisWorker{
		util.GetRedisPool(options.Connection),
		NewWorker(),
	}
	return
}

func (w *RedisWorker) Process(options *proto.RedisWorkerProcessOptions) {
	w.start.Wait()

	log.WithFields(log.Fields{
		"concurrency": options.Concurrency,
		"queue":       options.Queue,
	}).Info("Initiating workers")
	for i := 0; i < options.Concurrency; i++ {
		go w.process(options, i)
	}
}

//
// func (w *RedisWorker) process(options *proto.RedisWorkerProcessOptions, id int) {
// 	conn := w.pool.Get()
// 	defer conn.Close()
//
// 	log.WithFields(log.Fields{
// 		"id":    id,
// 		"queue": options.Queue,
// 	}).Info("Worker thread connected to Redis")
//
// 	for {
// 		message, err := conn.Do("BRPOP", options.Queue, 0)
// 		if err != nil {
// 			log.WithFields(log.Fields{
// 				"id":    id,
// 				"error": err,
// 				"queue": options.Queue,
// 			}).Error("Redis BRPOP error")
// 			return
// 		}
// 		log.WithFields(log.Fields{
// 			"id":    id,
// 			"queue": options.Queue,
// 		}).Info("Redis BRPOP success")
// 		// FIXME: check length, and if it's bytes
// 		switch m := message.(type) {
// 		case []bytes:
// 			buff := bytes.NewBuffer(m)
// 			msg, err := encoders.GobDecode(buff)
// 			if err != nil {
// 				log.WithFields(log.Fields{
// 					"id":      id,
// 					"message": m,
// 					"error":   err,
// 					"queue":   options.Queue,
// 				}).Error("Error decoding message")
// 				continue
// 			}
// 			processMessage(msg)
// 		default:
// 			log.WithFields(log.Fields{
// 				"id":      id,
// 				"message": m,
// 				"type":    reflect.TypeOf(message),
// 				"queue":   options.Queue,
// 			}).Error("Invalid message type received")
// 			continue
// 		}
// 	}
// }

func (w *RedisWorker) process(options *proto.RedisWorkerProcessOptions, id int) {
	conn := w.pool.Get()
	defer conn.Close()

	log.WithFields(log.Fields{
		"id":    id,
		"queue": options.Queue,
	}).Info("Worker thread connected to Redis")

	for {
		message, err := conn.Do("BRPOP", options.Queue, 0)
		if err != nil {
			log.WithFields(log.Fields{
				"id":    id,
				"error": err,
				"queue": options.Queue,
			}).Error("Redis BRPOP error")
			return
		}
		log.WithFields(log.Fields{
			"id":    id,
			"queue": options.Queue,
		}).Debug("Redis BRPOP success")
		// FIXME: check length, and if it's bytes
		switch m := message.(type) {
		case []interface{}:
			switch mbytes := m[1].(type) {
			case []byte:
				buff := bytes.NewBuffer(mbytes)
				msg, err := encoders.GobDecode(buff)
				if err != nil {
					log.WithFields(log.Fields{
						"id":      id,
						"message": mbytes,
						"error":   err,
						"queue":   options.Queue,
					}).Error("Error decoding message")
					continue
				}
				result := processMessage(msg)
				log.WithFields(log.Fields{
					"WorkerID":     id,
					"JobID":        msg.JobID,
					"Job.FuncName": msg.Job.FuncName,
					"Job.FuncArgs": msg.Job.FuncArgs,
					"Result":       result,
				}).Info("Finished job")

				// var outbuff bytes.Buffer
				// msg, err := encoders.GobEncode(outbuff)

				// err = conn.Do("LPUSH", result)
			default:
				log.WithFields(log.Fields{
					"id":      id,
					"message": m,
					"type":    reflect.TypeOf(message),
					"queue":   options.Queue,
				}).Error("Invalid message type received, take 2")
			}
		default:
			log.WithFields(log.Fields{
				"id":      id,
				"message": m,
				"type":    reflect.TypeOf(message),
				"queue":   options.Queue,
			}).Error("Invalid message type received")
			continue
		}
	}
}
