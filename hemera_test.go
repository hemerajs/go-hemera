package hemera

import (
	"errors"
	"fmt"
	"testing"
	"time"

	natsServer "github.com/nats-io/gnatsd/server"
	gnatsd "github.com/nats-io/gnatsd/test"
	nats "github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

const TEST_PORT = 8368

var reconnectOpts = nats.Options{
	Url:            fmt.Sprintf("nats://localhost:%d", TEST_PORT),
	AllowReconnect: true,
	MaxReconnect:   10,
	ReconnectWait:  100 * time.Millisecond,
	Timeout:        nats.DefaultTimeout,
}

// Dumb wait program to sync on callbacks, etc... Will timeout
func Wait(ch chan bool) error {
	return WaitTime(ch, 5*time.Second)
}

func WaitTime(ch chan bool, timeout time.Duration) error {
	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
	}
	return errors.New("timeout")
}

func RunServerOnPort(port int) *natsServer.Server {
	opts := gnatsd.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(opts)
}

func RunServerWithOptions(opts natsServer.Options) *natsServer.Server {
	return gnatsd.RunServer(&opts)
}

type MathPattern struct {
	Topic string `json:"topic"`
	Cmd   string `json:"cmd"`
}

type RequestPattern struct {
	Topic string `json:"topic" mapstructure:"topic"`
	Cmd   string `json:"cmd" mapstructure:"cmd"`
	A     int    `json:"a" mapstructure:"a"`
	B     int    `json:"b" mapstructure:"b"`
}

type Response struct {
	Result int `json:"result"`
}

func TestCreateHemera(t *testing.T) {
	assert := assert.New(t)

	ts := RunServerOnPort(TEST_PORT)
	defer ts.Shutdown()

	opts := reconnectOpts
	nc, err := opts.Connect()
	defer nc.Close()

	if err != nil {
		panic(err)
	}

	nc.Flush()

	h, _ := CreateHemera(nc)

	assert.NotEqual(h, nil, "they should not nil")

}

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	ts := RunServerOnPort(TEST_PORT)
	defer ts.Shutdown()

	opts := reconnectOpts
	nc, err := opts.Connect()
	defer nc.Close()

	if err != nil {
		panic(err)
	}

	h, _ := CreateHemera(nc)

	pattern := MathPattern{Topic: "math", Cmd: "add"}

	h.Add(pattern, func(req *RequestPattern, reply Reply, context Context) {
		reply.Send(Response{Result: req.A + req.B})
	})

	nc.Flush()

	assert.Equal(h.Router.Len(), 1, "Should be 1")

}

func TestActRequest(t *testing.T) {
	assert := assert.New(t)
	ch := make(chan bool)
	actResult := &Response{}

	ts := RunServerOnPort(TEST_PORT)
	defer ts.Shutdown()

	opts := reconnectOpts
	nc, err := opts.Connect()
	defer nc.Close()

	if err != nil {
		panic(err)
	}

	h, _ := CreateHemera(nc)

	pattern := MathPattern{Topic: "math", Cmd: "add"}

	h.Add(pattern, func(req *RequestPattern, reply Reply, context Context) {
		reply.Send(Response{Result: req.A + req.B})
	})

	nc.Flush()

	requestPattern := RequestPattern{Topic: "math", Cmd: "add", A: 1, B: 2}
	go h.Act(requestPattern, func(resp *Response, err Error, context Context) {
		actResult = resp
		ch <- true
	})

	nc.Flush()

	if err := Wait(ch); err != nil {
		t.Fatal("Did not receive our message")
	}

	assert.Equal(actResult.Result, 3, "Should be 3")

}

func TestNoDuplicatesAllowed(t *testing.T) {
	assert := assert.New(t)

	ts := RunServerOnPort(TEST_PORT)
	defer ts.Shutdown()

	opts := reconnectOpts
	nc, err := opts.Connect()
	defer nc.Close()

	if err != nil {
		panic(err)
	}

	h, _ := CreateHemera(nc)

	pattern := MathPattern{Topic: "math", Cmd: "add"}

	h.Add(pattern, func(req *RequestPattern, reply Reply, context Context) {
		reply.Send(Response{Result: req.A + req.B})
	})

	_, errAdd := h.Add(pattern, func(req *RequestPattern, reply Reply, context Context) {
		reply.Send(Response{Result: req.A + req.B})
	})

	nc.Flush()

	assert.Equal(errAdd.Error(), "Pattern is already registered", "Should be not allowed to add duplicate patterns")

}
