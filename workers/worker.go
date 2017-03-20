package workers

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/proto"
	"github.com/xLegoz/go-swallow/util"
)

var access sync.Mutex
var started bool

type Worker struct {
	exit  *sync.WaitGroup
	start *sync.WaitGroup
}

func NewWorker() *Worker {
	// Blocks start from the beginning
	start := &sync.WaitGroup{}
	start.Add(1)
	return &Worker{
		exit:  &sync.WaitGroup{},
		start: start,
	}
}

func (w *Worker) Run() {
	log.Debug("Worker run")
	w.Start()
	go util.HandleSignals(w.Quit)
	w.waitForExit()
}

func (w *Worker) Start() {
	log.Info("Worker start")
	access.Lock()
	defer access.Unlock()
	defer w.start.Done()
	defer w.exit.Add(1)

	if started {
		log.Info("Worker already started")
		return
	}

	started = true
}

func (w *Worker) Quit() {
	log.Info("Worker quit")
	access.Lock()
	defer access.Unlock()

	if !started {
		log.Info("Worker already stopped")
		return
	}

	w.exit.Done()
	w.waitForExit()

	started = false
}

func (w *Worker) waitForExit() {
	w.exit.Wait()
}

func processMessage(msg *proto.Message) interface{} {
	log.WithFields(log.Fields{
		"jobID": msg.JobID,
		"job":   msg.Job,
	}).Debug("Calling job for message")
	return util.CallWith(msg.Job)
}
