package hemera

import (
	"testing"

	server "github.com/nats-io/gnatsd/server"
	gnatsd "github.com/nats-io/gnatsd/test"
	nats "github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

const testPort = 8368

func RunServerOnPort(port int) *server.Server {
	opts := gnatsd.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(opts)
}

func RunServerWithOptions(opts server.Options) *server.Server {
	return gnatsd.RunServer(&opts)
}

func CreateHemera(t *testing.T) {
	assert := assert.New(t)

	ts := RunServerOnPort(testPort)
	nc, _ := nats.Connect(nats.DefaultURL)
	hr, _ := NewHemera(nc)

	assert.NotEqual(hr, nil, "they should not nil")

	ts.Shutdown()

}
