/*
WIP Client
*/

package clients

import (
	"fmt"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/encoders"
	"github.com/xLegoz/go-swallow/proto"
	"github.com/xLegoz/go-swallow/util"
)

type DirectClient struct {
	options    *proto.DirectClientOptions
	connection net.Conn
}

func NewDirectClient(options *proto.DirectClientOptions) (directClient *DirectClient, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", options.Address, options.Port))
	if err != nil {
		log.WithFields(log.Fields{
			"address": options.Address,
			"port":    options.Port,
			"error":   err,
		}).Error("Direct connect failed")
		return nil, err
	}
	directClient = &DirectClient{
		options:    options,
		connection: conn,
	}
	return
}

func (c *DirectClient) Perform(fn interface{}, args interface{}, extra ...interface{}) error {
	jobRequest, err := util.CreateJobRequest(fn, args)
	if err != nil {
		panic(err)
	}
	msg := &proto.Message{
		JobID:  util.GenerateJobID(),
		Job:    jobRequest,
		Status: 1,
	}
	err = encoders.GobEncode(c.connection, msg)
	if err != nil {
		panic(err)
	}

	// XXX: Placeholder, serialize and send message

	return nil
}

func (c *DirectClient) Close() {
	c.connection.Close()
}
