package proto

type RedisWorkerProcessOptions struct {
	Queue       string
	Concurrency int
}

type RedisWorkerOptions struct {
	Connection *RedisConnectionOptions
}

type RedisConnectionOptions struct {
	Address  string
	Password string
	Database string
	Poolsize int
}

type RedisClientOptions struct {
	Connection *RedisConnectionOptions
	Queue      string
}

func (c *RedisConnectionOptions) SetDefaults() *RedisConnectionOptions {
	c.Address = "localhost:6379"
	c.Poolsize = 2
	return c
}

var DefaultRedisConnectionOptions = (&RedisConnectionOptions{}).SetDefaults()
