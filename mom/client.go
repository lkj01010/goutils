package mom

type Client interface {

	// Pub Sub
	Subscribe(sbuj string, h MsgHandler)
	UnSubscribe(subj string)
	Publish(subj string, v interface{}) error

	// Request Reply
	Request(subj string, v interface{}, vPtr interface{}) error

	// Message queue ...
}

type MsgHandler func(replySbuj string, msg []byte)