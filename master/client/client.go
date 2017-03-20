package client

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/lib"
	"github.com/xLegoz/go-swallow/master/worker"
)

/*
Client represents a connected client
*/
type Client struct {
	ID         int
	Job        *lib.JobRequest
	Connection net.Conn
}

var clients = []*Client{}
var clientsLock sync.Mutex

/*
Listen will listen for clients on port
*/
func Listen(port int) {
	log.WithFields(log.Fields{
		"port": port,
	}).Info("Listening for clients")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error listening for clients")
		return
	}
	clientNum := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error accepting connection for client")
			continue
		}
		go clientHandleConnection(conn, clientNum)
		clientNum++
	}
}

func clientHandleConnection(conn net.Conn, id int) {
	defer conn.Close()
	log.WithFields(log.Fields{
		"clientId":      id,
		"localAddress":  conn.LocalAddr(),
		"remoteAddress": conn.RemoteAddr(),
	}).Info("Client connected")

	// FIXME: how do we make this last more than one item?
	dec := gob.NewDecoder(conn)
	job := &lib.JobRequest{}
	dec.Decode(job)
	err := lib.CheckJobValidity(job)
	log.WithFields(log.Fields{
		"clientId": id,
		"job":      job,
	}).Info("Received job")
	if err != nil {
		log.WithFields(log.Fields{
			"clientId": id,
			"job":      job,
			"error":    err,
		}).Error("Job is invalid")
		return
	}

	go worker.SendToWorker(job)
}
