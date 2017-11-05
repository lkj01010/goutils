package mom
import (
	"time"
	"fmt"
	"github.com/nats-io/nats"
)


type natsImpl struct {
	nc            *nats.EncodedConn
	subscriptions map[string]*nats.Subscription
}

func NewNats() Client {
	return &natsImpl{
		nc: newDefaultJsonConn(),
		subscriptions: map[string]*nats.Subscription{},
	}
}

func (ni *natsImpl)Subscribe(subj string, h MsgHandler) {
	subscription, err := ni.nc.Subscribe(subj, func(msg *nats.Msg) {
//		log.Debugf("Subscribe subj=%+v, natsReply=%+v", subj, msg)
		h(msg.Reply, msg.Data);
	})
	if err != nil {
		panic("Subscribe " + subj + " err=" + err.Error());
	}

	ni.subscriptions[subj] = subscription;
}

func (ni *natsImpl)UnSubscribe(subj string) {
	if err := ni.subscriptions[subj].Unsubscribe(); err != nil {
		panic("Unsubscribe " + subj + err.Error());
	}
}

func (ni *natsImpl)Publish(subj string, v interface{}) error {
	return ni.nc.Publish(subj, v);
}

func (ni *natsImpl)Request(subj string, v interface{}, vPtr interface{}) error {
	return ni.nc.Request(subj, v, vPtr, 5000*time.Millisecond)
}

//////////////////////////////////////////////////

var defaultOpts = nats.Options{
	Url:            fmt.Sprintf("nats://localhost:%d", nats.DefaultPort),
	AllowReconnect: true,
	MaxReconnect:   10,
	ReconnectWait:  100 * time.Millisecond,
	Timeout:        nats.DefaultTimeout,
}

func newJsonEncodedConn(nc *nats.Conn) *nats.EncodedConn {
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic("Failed to create an encoded connection: err=" + err.Error())
	}
	return ec
}

func newDefaultJsonConn() *nats.EncodedConn {
	nc, err := defaultOpts.Connect()
	if err != nil {
		panic("Failed to connect gnatsd, err=" + err.Error())
	}
	ec := newJsonEncodedConn(nc)
	return ec
}
