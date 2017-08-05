package hemera

import (
	"encoding/json"

	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/nuid"
)

type Reply struct {
	Pattern interface{}
	Conn    *nats.Conn
	Reply   string
}

func (r *Reply) Send(payload interface{}) {
	response := packet{
		Pattern: r.Pattern,
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

	// Encode to JSON
	data, _ := json.Marshal(&response)

	// Send
	r.Conn.Publish(r.Reply, data)
}
