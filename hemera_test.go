package hemera

import (
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

func RunServerOnPort(port int) *natsServer.Server {
	opts := gnatsd.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(opts)
}

func RunServerWithOptions(opts natsServer.Options) *natsServer.Server {
	return gnatsd.RunServer(&opts)
}

type MathPattern struct {
	Topic string
	Cmd   string
}

type RequestPattern struct {
	Topic string
	Cmd   string
	A     int
	B     int
}

type Response struct {
	Result int
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

	assert.Equal(len(h.Router.List()), 1, "Should be 1")

}

func TestActRequest(t *testing.T) {
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

	requestPattern := RequestPattern{Topic: "math", Cmd: "add", A: 1, B: 2}
	res := &Response{}
	h.Act(requestPattern, res)

	assert.Equal(t, res.Result, 3, "Should be 3")

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

	assert.Equal(errAdd.Message, "Pattern is already registered", "Should be not allowed to add duplicate patterns")

}
