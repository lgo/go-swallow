package clients

type Client interface {
	Perform(fn interface{}, args interface{}, extra ...interface{}) (chan interface{}, error)
}
