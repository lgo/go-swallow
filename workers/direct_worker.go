package workers

import (
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/encoders"
	"github.com/xLegoz/go-swallow/proto"
	"github.com/xLegoz/go-swallow/util"
)

type DirectWorkerOptions struct {
	Port int
}

type DirectWorker struct {
	listener net.Listener
	*Worker
}

func NewDirectWorker(options *proto.DirectWorkerOptions) (directWorker *DirectWorker, err error) {
	log.WithFields(log.Fields{
		"options": options,
	}).Info("New Direct worker")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", options.Port))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error listening for clients")
		return nil, err
	}

	directWorker = &DirectWorker{
		ln,
		NewWorker(),
	}
	return
}

func (w *DirectWorker) Process() {
	clientNum := 0
	for {
		conn, err := w.listener.Accept()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error accepting connection for client")
			continue
		}
		go w.handleConnection(conn, clientNum)
		clientNum++
	}
}

func (w *DirectWorker) handleConnection(conn net.Conn, id int) {
	defer conn.Close()
	log.WithFields(log.Fields{
		"cliendID": id,
	}).Info("New client connected")
	for {
		msg, err := encoders.GobDecode(conn)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Decoding error")
			return
		}
		log.WithFields(log.Fields{
			"clientId": id,
			"job":      msg,
		}).Info("Received message")
		err = util.CheckJobValidity(msg.Job)
		if err != nil {
			log.WithFields(log.Fields{
				"clientId": id,
				"message":  msg,
				"error":    err,
			}).Error("Message is invalid")
			return
		}
		// FIXME: go routine this?
		util.CallWith(msg.Job)
	}
}
