package hemera

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	nats "github.com/nats-io/go-nats"
	"github.com/nats-io/nuid"
)

const (
	// RequestType represent the request with default request / reply semantic
	RequestType = "request"
	// PubsubType represent the request with publish / subscribe semantic
	PubsubType = "pubsub"
	// RequestTimeout is the maxiumum act timeout in miliseconds
	RequestTimeout = 2000
)

var (
	ErrAddTopicRequired           = errors.New("Topic is required")
	ErrActTopicRequired           = errors.New("Topic is required")
	ErrInvalidTopicType           = errors.New("Topic must be from type string")
	ErrInvalidMapping             = errors.New("Map could not be mapped to struct")
	ErrInvalidAddHandlerArguments = errors.New("Add Handler requires at least one argument")
)

func GetDefaultOptions() Options {
	opts := Options{
		Timeout: RequestTimeout,
	}
	return opts
}

// Option is a function on the options for hemera
type Option func(*Options) error

type Options struct {
	Timeout time.Duration
}

type actHandler func(ClientResult)
type Handler interface{}

// Hemera is the main struct
type Hemera struct {
	Conn   *nats.Conn
	Router Router
	Opts   Options
}

type request struct {
	ID          string `json:"id"`
	RequestType string `json:"type"`
}

type ClientResult interface{}

// Error is the default error struct
type Error struct {
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
	Pattern  interface{}            `json:"pattern"`
	Meta     map[string]interface{} `json:"meta"`
	Delegate map[string]interface{} `json:"delegate"`
	Result   interface{}            `json:"result"`
	Trace    trace                  `json:"trace"`
	Request  request                `json:"request"`
	Error    *Error                 `json:"error"`
}

// New create a new Hemera struct
func NewHemera(conn *nats.Conn, options ...Option) (Hemera, error) {
	opts := GetDefaultOptions()
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return Hemera{Opts: opts, Router: Router{}}, err
		}
	}
	return Hemera{Conn: conn, Opts: opts, Router: Router{}}, nil
}

// Timeout is an Option to set the timeout for a act request
func Timeout(t time.Duration) Option {
	return func(o *Options) error {
		o.Timeout = t
		return nil
	}
}

// Add is a method to subscribe on a specific topic
func (h *Hemera) Add(p interface{}, cb Handler) (*nats.Subscription, error) {
	s := structs.New(p)
	f := s.Field("Topic")

	if f.IsZero() {
		return nil, ErrAddTopicRequired
	}

	topic, ok := f.Value().(string)

	if !ok {
		return nil, ErrInvalidTopicType
	}

	// Get the types of the Add handler args
	argMsgType, argReplyType, numArgs := argInfo(cb)

	if numArgs != 2 || argMsgType == nil || argReplyType == nil {
		return nil, ErrInvalidAddHandlerArguments
	}

	cbValue := reflect.ValueOf(cb)

	natsCB := func(m *nats.Msg) {
		var oPtr reflect.Value
		if argMsgType.Kind() != reflect.Ptr {
			oPtr = reflect.New(argMsgType)
		} else {
			oPtr = reflect.New(argMsgType.Elem())
		}

		// Get "Value" of the reply callback for the reflection Call
		reply := Reply{Pattern: p, Conn: h.Conn, Reply: m.Reply}

		oReplyPtr := reflect.ValueOf(reply)

		pack := packet{}

		// decoding hemera packet
		json.Unmarshal(m.Data, &pack)

		// Pattern is the request
		o := pack.Pattern

		// return the value of oPtr as interface {}
		oi := oPtr.Interface()

		// Decode map to struct
		err := mapstructure.Decode(o, oi)

		if err != nil {
			panic(err)
		}

		// Get "Value" of the reply callback for the reflection Call
		oPtr = reflect.ValueOf(oi)

		// array of arguments for the callback handler
		oV := []reflect.Value{oPtr, oReplyPtr}

		cbValue.Call(oV)
	}

	return h.Conn.QueueSubscribe(topic, topic, natsCB)
}

// Act is a method to send a message to a NATS subscriber which the specific topic
func (h *Hemera) Act(p interface{}, handler actHandler) (bool, error) {

	s := structs.New(p)
	f := s.Field("Topic")

	if f.IsZero() {
		return false, ErrActTopicRequired
	}

	topic, ok := f.Value().(string)

	if !ok {
		return false, ErrInvalidTopicType
	}

	request := packet{
		Pattern: p,
		Request: request{
			ID:          nuid.Next(),
			RequestType: RequestType,
		},
	}

	data, _ := json.Marshal(&request)
	m, err := h.Conn.Request(topic, data, h.Opts.Timeout*time.Millisecond)

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	pack := packet{}
	mErr := json.Unmarshal(m.Data, &pack)

	if mErr != nil {
		log.Fatal(mErr)
		return false, err
	}

	if pack.Error != nil {
		handler(pack.Error)
	} else {
		handler(pack.Result)
	}

	return true, nil
}

// Dissect the cb Handler's signature
func argInfo(cb Handler) (reflect.Type, reflect.Type, int) {
	cbType := reflect.TypeOf(cb)

	if cbType.Kind() != reflect.Func {
		panic("nats: Handler needs to be a func")
	}

	numArgs := cbType.NumIn()

	if numArgs < 2 {
		return nil, nil, numArgs
	}

	return cbType.In(0), cbType.In(1), numArgs
}
