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
	RequestTimeout    = 2000
	DepthIndexing     = true
	InsertionIndexing = false
)

var (
	ErrAddTopicRequired           = errors.New("Topic is required")
	ErrActTopicRequired           = errors.New("Topic is required")
	ErrInvalidTopicType           = errors.New("Topic must be from type string")
	ErrInvalidMapping             = errors.New("Map could not be mapped to struct")
	ErrInvalidAddHandlerArguments = errors.New("Add Handler requires at least two argument")
	ErrInvalidActHandlerArguments = errors.New("Act Handler requires at least two argument")
	ErrPatternNotFound            = errors.New("Pattern not found")
	ErrDuplicatePattern           = errors.New("Pattern is already registered")
)

func GetDefaultOptions() Options {
	opts := Options{
		Timeout:          RequestTimeout,
		IndexingStrategy: false,
	}
	return opts
}

// Option is a function on the options for hemera
type Option func(*Options) error

type Options struct {
	Timeout          time.Duration
	IndexingStrategy bool
}

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
	Pattern  interface{} `json:"pattern"`
	Meta     interface{} `json:"meta"`
	Delegate interface{} `json:"delegate"`
	Result   interface{} `json:"result"`
	Trace    trace       `json:"trace"`
	Request  request     `json:"request"`
	Error    *Error      `json:"error"`
}

// New create a new Hemera struct
func CreateHemera(conn *nats.Conn, options ...Option) (Hemera, error) {
	opts := GetDefaultOptions()
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return Hemera{Opts: opts, Router: NewRouter(opts.IndexingStrategy)}, err
		}
	}
	return Hemera{Conn: conn, Opts: opts, Router: NewRouter(opts.IndexingStrategy)}, nil
}

// Timeout is an Option to set the timeout for a act request
func Timeout(t time.Duration) Option {
	return func(o *Options) error {
		o.Timeout = t
		return nil
	}
}

func IndexingStrategy(isDeep bool) Option {
	return func(o *Options) error {
		o.IndexingStrategy = isDeep
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
	argTypes, numArgs := ArgInfo(cb)

	if numArgs < 2 {
		return nil, ErrInvalidAddHandlerArguments
	}

	lp := h.Router.Lookup(p)

	if lp != nil {
		return nil, ErrDuplicatePattern
	}

	h.Router.Add(p, cb)

	// Response struct
	argMsgType := argTypes[0]

	return h.Conn.QueueSubscribe(topic, topic, func(m *nats.Msg) {
		h.callAddAction(topic, m, argMsgType, numArgs)
	})
}

func (h *Hemera) callAddAction(topic string, m *nats.Msg, mContainer reflect.Type, numArgs int) {
	var oPtr reflect.Value
	if mContainer.Kind() != reflect.Ptr {
		oPtr = reflect.New(mContainer)
	} else {
		oPtr = reflect.New(mContainer.Elem())
	}

	pack := packet{}

	// decoding hemera packet
	json.Unmarshal(m.Data, &pack)

	context := Context{Meta: pack.Meta, Delegate: pack.Delegate, Trace: pack.Trace}

	oContextPtr := reflect.ValueOf(context)

	// Pattern is the request
	o := pack.Pattern

	// return the value of oPtr as interface {}
	oi := oPtr.Interface()

	// Decode map to struct
	err := mapstructure.Decode(o, oi)

	if err != nil {
		panic(err)
	}

	e := oPtr.Elem().Interface()

	p := h.Router.Lookup(e)

	if p != nil {
		// Get "Value" of the reply callback for the reflection Call
		reply := Reply{Pattern: p.Pattern, Reply: m.Reply, Hemera: h}

		oReplyPtr := reflect.ValueOf(reply)

		cbValue := reflect.ValueOf(p.Payload)

		// Get "Value" of the reply callback for the reflection Call

		oPtr = reflect.ValueOf(oi)

		// array of arguments for the callback handler
		var oV []reflect.Value

		if numArgs == 2 {
			oV = []reflect.Value{oPtr, oReplyPtr}
		} else {
			oV = []reflect.Value{oPtr, oReplyPtr, oContextPtr}
		}

		cbValue.Call(oV)
	} else {
		log.Fatal(ErrPatternNotFound)
	}
}

// Act is a method to send a message to a NATS subscriber which the specific topic
func (h *Hemera) Act(p interface{}, cb Handler) (bool, error) {

	s := structs.New(p)
	topicField := s.Field("Topic")

	if topicField.IsZero() {
		return false, ErrActTopicRequired
	}

	topic, ok := topicField.Value().(string)

	if !ok {
		return false, ErrInvalidTopicType
	}

	var metaField interface{}
	if field, ok := s.FieldOk("Meta_"); ok {
		metaField = field.Value()
	}

	var delegateField interface{}
	if field, ok := s.FieldOk("Delegate_"); ok {
		delegateField = field.Value()
	}

	pattern := CleanPattern(p)

	argTypes, numArgs := ArgInfo(cb)

	if numArgs < 2 {
		return false, ErrInvalidActHandlerArguments
	}

	// Response struct
	argMsgType := argTypes[0]

	cbValue := reflect.ValueOf(cb)

	var oPtr reflect.Value
	if argMsgType.Kind() != reflect.Ptr {
		oPtr = reflect.New(argMsgType)
	} else {
		oPtr = reflect.New(argMsgType.Elem())
	}

	request := packet{
		Pattern:  pattern,
		Meta:     metaField,
		Delegate: delegateField,
		Trace: trace{
			TraceID: nuid.Next(),
		},
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

	// return the value of oPtr as interface {}
	oi := oPtr.Interface()

	// Pattern is the request
	o := pack.Result

	// Decode result map to struct
	errResultMap := mapstructure.Decode(o, oi)

	if errResultMap != nil {
		panic(errResultMap)
	}

	// Get "Value" of the reply callback for the reflection Call
	oPtr = reflect.ValueOf(oi)

	errMsg := pack.Error

	// create container for error
	errorMsg := Error{}

	if errMsg != nil {
		// Decode error map to struct
		errErrMap := mapstructure.Decode(errMsg, &errorMsg)

		if errErrMap != nil {
			panic(errErrMap)
		}
	}

	context := Context{Meta: pack.Meta, Delegate: pack.Delegate, Trace: pack.Trace}

	oContextPtr := reflect.ValueOf(context)

	if pack.Error != nil {
		errVal := reflect.ValueOf(errorMsg)

		var argValues []reflect.Value
		if numArgs == 2 {
			argValues = []reflect.Value{oPtr, errVal}
		} else {
			argValues = []reflect.Value{oPtr, errVal, oContextPtr}
		}
		cbValue.Call(argValues)
	} else {
		errVal := reflect.ValueOf(errorMsg)

		var argValues []reflect.Value
		if numArgs == 2 {
			argValues = []reflect.Value{oPtr, errVal}
		} else {
			argValues = []reflect.Value{oPtr, errVal, oContextPtr}
		}
		cbValue.Call(argValues)
	}

	return true, nil
}
