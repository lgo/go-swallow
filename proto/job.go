package proto

import "encoding/gob"

/*
JobRequest is a representitive of an over-the-wire Job
- JobID is used in the response of this job to help identify it
- FuncName is to denote the function in the registry for the job
- FuncArgs are the arguments to call with
*/
type JobRequest struct {
	JobID    string
	FuncName string
	FuncArgs interface{}
}

/*
JobResponse represents the over-the-wire job response
- JobID is used in the response to associate it with a request
- FuncReturn are the return values of the function call
*/
type JobResponse struct {
	JobID      string
	FuncReturn []interface{}
}

type Message struct {
	JobID  string
	Job    *JobRequest
	Status int // ???
}

type EnqueueOptions struct {
	Retry bool
	At    int
}

func init() {
	gob.Register(new(JobRequest))
	gob.Register(new(JobResponse))
	gob.Register(new(Message))
}
