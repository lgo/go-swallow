# go-swallow

go-swallow is a Go job library inspired from the ease of usability provided by Sidekiq.

Using go-swallow is as as simple as starting up your worker with `worker.Run` and queue jobs with `client.Perform`!

```go
// Create my job
func Add(args *Args) {
  result := args.A + args.B
  log.Infof("Look ma, I added! %d + %d = %d", args.A, args.B, result)
  return result
}

func serverMain() {
  // Initialize the job worker
  ...
  // Start processing jobs!
  w.Run()
}


func clientMain() {
  // Initialize the job queue connection
  ...
  // Start a job!
  client.Perform(jobs.Add, jobs.Args{A: 1, B: 2})
}
```

(see [example usage](https://github.com/xLegoz/go-swallow#example-usage) for more details, including initialization)

### Features
- Worker and client Redis backend
- Registering functions as jobs
- Ability to specify concurrency, and have multiple workers

### Future features
- Reliability guarantees (i.e. use `BRPOPLPUSH` to keep persist jobs)
- Implement additional backends
  - Direct TCP
  - PostgreSQL
- Allow dynamic perform arguments (will have unvalidated types)
- Provide a channel for the client to retrieve results
- Safely recover from worker errors, and passing the error on to client channels

### Issues
- Registering duplicate or primitive types as arguments or return signatures

## Example usage

See `examples/redis` for an example client and worker. Below is a shorter copy of the three files needed to get started.

### Jobs
You need to create your job functions and register them with go-swallow. Jobs must have a single `struct` argument, which must be registered with `go-swallow/util.Register`

`jobs.go`
```go
package jobs

import (
	"log"
	"github.com/xLegoz/go-swallow/util"
)

type Args struct {
	A int
	B int
}

func Add(args *Args) {
	result := args.A + args.B
	log.Infof("Look ma, I added! %d + %d = %d", args.A, args.B, result)
	return result
}

func init() {
	util.Register(Add, new(Args))
}
```

### Client
Your client is implemented using the client library, and can be integrated into regular Go code. Creating a job is as simple as calling `Perform` on your desired function!

`client.go`
```go
package main

import (
	"os"

	_ "jobs"
	"github.com/xLegoz/go-swallow/clients"
	"github.com/xLegoz/go-swallow/proto"
)

func main() {
	client, _ := clients.NewRedisClient(&proto.RedisClientOptions{
		Queue: "go-swallow-jobs",
	})

	client.Perform(jobs.Add, jobs.Args{A: 1, B: 2})
}
```

### Worker
The worker is largely boilerplate, and doesn't need to be changed to accommodate new jobs. It simply needs to import your jobs, set the worker to process a Redis queue, then `Run`!

`worker.go`
```go
package main

import (
	"os"

	_ "jobs"
	"github.com/xLegoz/go-swallow/workers"
	"github.com/xLegoz/go-swallow/proto"
)

func main() {
	w, _ := workers.NewRedisWorker(&proto.RedisWorkerOptions{})

	go client.Process(&proto.RedisWorkerProcessOptions{
		Queue:       "go-swallow-jobs",
		Concurrency: 3,
	})

	w.Run()
}
```
