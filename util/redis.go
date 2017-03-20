package util

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/xLegoz/go-swallow/proto"
)

func GetRedisPool(options *proto.RedisConnectionOptions) *redis.Pool {
	// FIXME: test that the connection is valid before hand, and return an error
	return &redis.Pool{
		MaxIdle:     options.Poolsize,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", options.Address)
			if err != nil {
				log.WithFields(log.Fields{
					"error":   err,
					"address": options.Address,
				}).Debug("Failed to connect to Redis")
				return nil, err
			}
			if options.Password != "" {
				if _, err = c.Do("AUTH", options.Password); err != nil {
					c.Close()
					log.WithFields(log.Fields{
						"error":    err,
						"address":  options.Address,
						"password": options.Password,
					}).Debug("Failed to authenticate Redis")
					return nil, err
				}
			}
			if options.Database != "" {
				if _, err = c.Do("SELECT", options.Database); err != nil {
					c.Close()
					log.WithFields(log.Fields{
						"error":    err,
						"address":  options.Address,
						"database": options.Database,
					}).Debug("Failed to select database in Redis")
					return nil, err
				}
			}
			log.WithFields(log.Fields{
				"address": options.Address,
			}).Debug("Connected to Redis")
			return c, err
		},
		TestOnBorrow: redisTestOnBorrow,
	}
}

func redisTestOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Debug("Redis pool connection PING failed")
	}
	return err
}
