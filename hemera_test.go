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

func CreateHemera(t *testing.T) {
	assert := assert.New(t)

	ts := RunServerOnPort(testPort)
	nc, _ := nats.Connect(nats.DefaultURL)
	h, _ := Create(nc, Timeout(2000))

	assert.NotEqual(h, nil, "they should not nil")

	ts.Shutdown()

}

func ActRequest(t *testing.T) {
	assert := assert.New(t)
	ch := make(chan bool)
	actResult := float64(0)

	type Result struct {
		Result float64 `json:"result"`
	}

	ts := RunServerOnPort(testPort)
	nc, _ := nats.Connect(nats.DefaultURL)
	h, _ := Create(nc)

	pattern := MathPattern{Topic: "math", Cmd: "add"}

	h.Add(pattern, func(context Context, req *RequestPattern, reply Reply) {
		reply.Send(req.A + req.B)
	})

	requestPattern := RequestPattern{Topic: "math", Cmd: "add", A: 1, B: 2}
	h.Act(requestPattern, func(context Context, err Error, resp *Response) {
		ch <- true
	})

	nc.Flush()

	ts.Shutdown()

	if err := Wait(ch); err != nil {
		assert.Equal(actResult, 3, "Should be 3")
	}

}
