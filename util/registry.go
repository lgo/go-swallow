package util

import (
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/xLegoz/go-swallow/proto"
)

// Error codes returned by failures in the Job pipeline
var (
	ErrJobNotRegistered         = errors.New("jobber: Job not registered")
	ErrJobIncorrectArugmentType = "jobber: Job arguments of incorrect type. Got %v, expected %v."
)

type registeredFunc struct {
	Func    interface{}
	ArgType interface{}
	/*ResultType interface{}*/
}

var registry = make(map[string]*registeredFunc)

func getFunctionName(i interface{}) (name string) {
	name = runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	log.WithFields(log.Fields{
		"name":      name,
		"interface": i,
	}).Debug("Get function name")
	return
}

/*
Register sets the function in the registry, along with the struct types
of it's argument and result
*/
func Register(fn interface{}, args interface{} /*, result interface{}*/) {
	// Register types with Gob decoder/encoder
	if args != nil {
		log.WithFields(log.Fields{
			"value": args,
		}).Debug("Gob registered args")
		gob.Register(args)
	}
	/*if result != nil {
		log.WithFields(log.Fields{
			"value": result,
		}).Debug("Gob registered result")
		gob.Register(result)
	}*/

	// Store the registered func
	registry[getFunctionName(fn)] = &registeredFunc{
		Func:    fn,
		ArgType: reflect.TypeOf(args),
		/*ResultType: reflect.TypeOf(result), */
	}

	log.WithFields(log.Fields{
		"name":    getFunctionName(fn),
		"argType": reflect.TypeOf(args),
		/*"resultType": reflect.TypeOf(result),*/
	}).Debug("Registered function")
}

/*
CreateJobRequest initializes a JobRequest, and validates it
*/
func CreateJobRequest(fn interface{}, args interface{}) (request *proto.JobRequest, err error) {
	// Create the job request
	request = &proto.JobRequest{
		FuncName: getFunctionName(fn),
		FuncArgs: args,
	}

	// Check validity
	err = CheckJobValidity(request)

	log.WithFields(log.Fields{
		"name": getFunctionName(fn),
		"args": args,
	}).Debug("Creating job request")
	return
}

/*
CheckJobValidity will ensure that the JobRequest has a valid
job function along with the correct argument type for the function
*/
func CheckJobValidity(jobRequest *proto.JobRequest) (err error) {
	log.WithFields(log.Fields{
		"jobRequest": jobRequest,
	}).Debug("Checking job validity")
	// Lookup job
	funcEntry := registry[jobRequest.FuncName]
	err = nil
	if funcEntry == nil {
		log.WithFields(log.Fields{
			"funcName": jobRequest.FuncName,
		}).Debug("Job func not registered")
		err = ErrJobNotRegistered
	} else if funcEntry.ArgType != reflect.TypeOf(jobRequest.FuncArgs) {
		// FIXME: Problems if the jobRequest.FuncArgs is not a pointer, unclear
		log.WithFields(log.Fields{
			"argType":      reflect.TypeOf(jobRequest.FuncArgs),
			"expectedType": funcEntry.ArgType,
		}).Debug("Job args invalid")
		err = fmt.Errorf(ErrJobIncorrectArugmentType, reflect.TypeOf(jobRequest.FuncArgs), funcEntry.ArgType)
	}
	return
}

/*
CallWith calls the function for a JobRequest
*/
func CallWith(jobRequest *proto.JobRequest) []interface{} {
	log.WithFields(log.Fields{
		"jobRequest": jobRequest,
	}).Debug("Calling job for JobRequest")
	registeredFunc := registry[jobRequest.FuncName]
	fn := reflect.ValueOf(registeredFunc.Func)
	results := fn.Call([]reflect.Value{reflect.ValueOf(jobRequest.FuncArgs)})
	returnResults := []interface{}{}
	for _, result := range results {
		returnResults = append(returnResults, result.Interface())
	}
	log.WithFields(log.Fields{
		"results": returnResults,
	}).Debug("Returning results for job")
	return returnResults
}
