package worker

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/lib"
)

// Status reflects the state of a connected worker
type Status int

// Statuses
const (
	Idle    Status = iota
	Running Status = iota
	Dead    Status = iota
)

/*
Worker represents a connected worker
*/
type Worker struct {
	ID            int
	Status        Status
	CurrentJob    *lib.JobRequest
	Connection    net.Conn
	ResultChannel chan int
}

var workers = []*Worker{}
var workersLock sync.Mutex
var workChannel = make(chan *lib.JobRequest)

/*
SendToWorker will enqueue a job
*/
func SendToWorker(job *lib.JobRequest) {
	workChannel <- job
}

/*
Listen will listen for workers on a port
*/
func Listen(port int) {
	log.WithFields(log.Fields{
		"port": port,
	}).Info("Listening for workers")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error listening for workers")
	}

	workerNum := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error accepting connection for worker")
			continue
		}
		go workerHandleConnection(conn, workerNum)
		workerNum++
	}
}

func addWorker(worker *Worker) {
	workersLock.Lock()
	defer workersLock.Unlock()
	workers = append(workers, worker)
}

func workerHandleConnection(conn net.Conn, id int) error {
	defer conn.Close()
	log.WithFields(log.Fields{
		"workerId":      id,
		"localAddress":  conn.LocalAddr(),
		"remoteAddress": conn.RemoteAddr(),
	}).Info("Worker connected")

	myWorker := &Worker{
		ID:            id,
		Status:        Idle,
		CurrentJob:    nil,
		Connection:    conn,
		ResultChannel: make(chan int),
	}

	addWorker(myWorker)

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	for {
		work := <-workChannel
		myWorker.Status = Running
		myWorker.CurrentJob = work

		log.WithFields(log.Fields{
			"workerId": myWorker.ID,
			"job":      myWorker.CurrentJob,
		}).Info("Sending work")

		enc.Encode(work)
		resp := &lib.JobResponse{}
		err := dec.Decode(resp)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Worker error")
			log.WithFields(log.Fields{
				"job": work,
			}).Warn("Re-queueing job")
			workChannel <- work
			return err
		}

		log.WithFields(log.Fields{
			"workerId": myWorker.ID,
			"resp":     resp,
		}).Info("Receving completed job")

		myWorker.Status = 0
	}
}
