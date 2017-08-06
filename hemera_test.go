package hemera

import (
	"errors"
	"testing"
	"time"

	natsServer "github.com/nats-io/gnatsd/server"
	gnatsd "github.com/nats-io/gnatsd/test"
	nats "github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

const testPort = 8368

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

	ts := RunServerOnPort(testPort)
	defer ts.Shutdown()

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		panic(err)
	}

	h, _ := Create(nc)

	assert.NotEqual(h, nil, "they should not nil")

}

func TestActRequest(t *testing.T) {
	assert := assert.New(t)
	ch := make(chan bool)
	actResult := &Response{}

	ts := RunServerOnPort(testPort)
	defer ts.Shutdown()

	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()

	h, _ := Create(nc)

	pattern := MathPattern{Topic: "math", Cmd: "add"}

	h.Add(pattern, func(req *RequestPattern, reply Reply, context Context) {
		reply.Send(Response{Result: req.A + req.B})
	})

	requestPattern := RequestPattern{Topic: "math", Cmd: "add", A: 1, B: 2}
	h.Act(requestPattern, func(resp *Response, err Error, context Context) {
		ch <- true
		actResult = resp
	})

	nc.Flush()

	if err := nc.LastError(); err != nil {
		panic(err)
	}

	if err := Wait(ch); err != nil {
		assert.Equal(actResult.Result, 3, "Should be 3")
	}

}
