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

func CreateHemera(t *testing.T) {
	assert := assert.New(t)

	ts := RunServerOnPort(testPort)
	nc, _ := nats.Connect(nats.DefaultURL)
	h, _ := NewHemera(nc)

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
	h, _ := NewHemera(nc)

	pattern := Pattern{"topic": "math", "cmd": "add"}
	h.Add(pattern, func(req Pattern, reply Reply) {
		r := req["a"].(float64) + req["b"].(float64)
		reply(Result{Result: r})
	})

	requestPattern := Pattern{"topic": "math", "cmd": "add", "a": 1, "b": 2}
	h.Act(requestPattern, func(resp ClientResult) {
		ch <- true
		actResult = resp.(float64)
	})

	nc.Flush()

	ts.Shutdown()

	if err := Wait(ch); err != nil {
		assert.Equal(actResult, 3, "Should be 3")
	}

}
