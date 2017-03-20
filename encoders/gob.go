package encoders

import (
	"encoding/gob"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/proto"
)

func GobEncode(writer io.Writer, msg *proto.Message) error {
	enc := gob.NewEncoder(writer)
	log.WithFields(log.Fields{
		"msg": msg,
	}).Debug("Encoding message")
	return enc.Encode(msg)
}

func GobDecode(reader io.Reader) (*proto.Message, error) {
	dec := gob.NewDecoder(reader)
	msg := &proto.Message{}
	err := dec.Decode(msg)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"msg": msg,
	}).Debug("Decoded message")
	return msg, nil
}
