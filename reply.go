package hemera

import (
	"github.com/json-iterator/go"
	"github.com/nats-io/nuid"
)

type Reply struct {
	Hemera  *Hemera
	Pattern interface{}
	Context *Context
	Reply   string
}

func (r *Reply) Send(payload interface{}) {
	response := packet{
		Pattern: r.Pattern,
		Meta:    r.Context.Meta,
		Trace:   r.Context.Trace,
		Request: request{
			ID:          nuid.Next(),
			RequestType: RequestType,
		},
	}

	// Check if error or message was passed
	he, ok := payload.(Error)
	if ok {
		response.Error = &he
	} else {
		response.Result = payload
	}

	data, _ := jsoniter.Marshal(&response)
	r.Hemera.Conn.Publish(r.Reply, data)
}
