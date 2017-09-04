package hemera

import (
	"github.com/json-iterator/go"
	"github.com/nats-io/nuid"
)

type Reply struct {
	hemera  *Hemera
	pattern interface{}
	context *Context
	reply   string
}

func (r *Reply) Send(payload interface{}) {
	response := packet{
		Pattern: r.pattern,
		Meta:    r.context.Meta,
		Trace:   r.context.Trace,
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
	r.hemera.Conn.Publish(r.reply, data)
}
