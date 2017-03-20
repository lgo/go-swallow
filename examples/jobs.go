package examples

import (
	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/util"
)

type ArithArgs struct {
	A int
	B int
}

func Add(args *ArithArgs) {
	result := args.A + args.B
	log.Infof("Add: %d + %d = %d", args.A, args.B, result)
}

func Test() {
	log.Info("Tested")
}

func init() {
	util.Register(Add, new(ArithArgs))
	util.Register(Test, nil)
}
