package server

import (
	"encoding/json"
	"fmt"

	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/nuid"
)

const (
	RequestType = "request"
	PubsubType  = "pubsub"
)

type Reply func(interface{})
type addHandler func(Request, Reply)
type actHandler func(error, Reply)
type Pattern map[string]interface{}

type Hemera struct {
	Conn *nats.Conn
}

type request struct {
	ID          string `json:"id"`
	RequestType string `json:"type"`
}

type HemeraError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int16  `json:"code"`
}

type trace struct {
	TraceID      string `json:"traceId"`
	ParentSpanID string `json:"parentSpanId"`
	SpanID       string `json:"spanId"`
	Timestamp    int64  `json:"timestamp"`
	Service      string `json:"service"`
	Method       string `json:"method"`
	Duration     int64  `json:"duration"`
}

type packet struct {
	Pattern  Pattern                `json:"pattern"`
	Meta     map[string]interface{} `json:"meta"`
	Delegate map[string]interface{} `json:"delegate"`
	Result   *interface{}           `json:"result"`
	Trace    trace                  `json:"trace"`
	Request  request                `json:"request"`
	Error    *HemeraError           `json:"error"`
}

func (h *Hemera) Add(p Pattern, handler addHandler) {
	topic := p["topic"].(string)
	h.Conn.Subscribe(topic, func(m *nats.Msg) {
		pack := packet{}
		json.Unmarshal(m.Data, &pack)
		fmt.Printf("Received a packet: %+v\n Reply: %s", pack, m.Reply)
		handler(Request{Payload: pack.Pattern}, func(payload interface{}) {
			response := packet{
				Pattern: p,
				Request: request{
					ID:          nuid.Next(),
					RequestType: RequestType,
				},
			}

			he, isError := payload.(HemeraError)
			if isError {
				response.Error = &he
			} else {
				response.Result = &payload
			}

			data, _ := json.Marshal(&response)
			fmt.Printf("\nPayload %s", string(data))
			h.Conn.Publish(m.Reply, data)
		})
	})
}

func (h *Hemera) Act(p Pattern, handler actHandler) {
}
