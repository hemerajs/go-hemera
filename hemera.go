package hemera

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/hemerajs/go-hemera/router"
	jsoniter "github.com/json-iterator/go"
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
	ErrAddTopicRequired           = "Topic is required"
	ErrActTopicRequired           = "Topic is required"
	ErrInvalidTopicType           = "Topic must be from type string"
	ErrInvalidMapping             = "Map could not be mapped to struct"
	ErrInvalidAddHandlerArguments = "Add Handler requires at least two argument"
	ErrInvalidActHandlerArguments = "Act Handler requires at least two argument"
	ErrPatternNotFound            = "Pattern not found"
	ErrDuplicatePattern           = "Pattern is already registered"
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
	Router *router.Router
	Opts   Options
}

type request struct {
	ID          string `json:"id"`
	RequestType string `json:"type"`
}

type Trace struct {
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
	Meta     Meta        `json:"meta"`
	Delegate Delegate    `json:"delegate"`
	Result   interface{} `json:"result"`
	Trace    Trace       `json:"trace"`
	Request  request     `json:"request"`
	Error    *Error      `json:"error"`
}

type Meta map[string]interface{}
type Delegate map[string]interface{}

// New create a new Hemera struct
func CreateHemera(conn *nats.Conn, options ...Option) (Hemera, error) {
	opts := GetDefaultOptions()
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return Hemera{Opts: opts, Router: router.NewRouter(opts.IndexingStrategy)}, err
		}
	}
	return Hemera{Conn: conn, Opts: opts, Router: router.NewRouter(opts.IndexingStrategy)}, nil
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
func (h *Hemera) Add(p interface{}, cb Handler) (*nats.Subscription, *Error) {
	s := structs.New(p)
	f := s.Field("Topic")

	if f.IsZero() {
		return nil, &Error{Message: ErrAddTopicRequired}
	}

	topic, ok := f.Value().(string)

	if !ok {
		return nil, &Error{Message: ErrInvalidTopicType}
	}

	// Get the types of the Add handler args
	argTypes, numArgs := ArgInfo(cb)

	if numArgs < 2 {
		return nil, &Error{Message: ErrInvalidAddHandlerArguments}
	}

	lp := h.Router.Lookup(p)

	if lp != nil {
		return nil, &Error{Message: ErrDuplicatePattern}
	}

	h.Router.Add(p, cb)

	// Response struct
	argMsgType := argTypes[0]

	sub, err := h.Conn.QueueSubscribe(topic, topic, func(m *nats.Msg) {
		h.callAddAction(topic, m, argMsgType, numArgs)
	})

	if err != nil {
		return nil, &Error{Message: err.Error()}
	}

	return sub, nil
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
	jsoniter.Unmarshal(m.Data, &pack)

	context := &Context{Trace: pack.Trace, Meta: pack.Meta, Delegate: pack.Delegate}

	oContextPtr := reflect.ValueOf(context)

	// Pattern is the request
	o := pack.Pattern

	// return the value of oPtr as an interface{}
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
		reply := Reply{
			Context: context,
			Pattern: p.Pattern,
			Reply:   m.Reply,
			Hemera:  h,
		}

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
func (h *Hemera) Act(p interface{}, out interface{}) *Context {

	context := &Context{}

	s := structs.New(p)
	topicField := s.Field("Topic")

	if topicField.IsZero() {
		context.Error = &Error{Message: ErrActTopicRequired}
		return context
	}

	topic, ok := topicField.Value().(string)

	if !ok {
		context.Error = &Error{Message: ErrInvalidTopicType}
		return context
	}

	var metaField Meta
	if field, ok := s.FieldOk("Meta"); ok {
		metaField = field.Value().(Meta)
	}

	var delegateField Delegate
	if field, ok := s.FieldOk("Delegate"); ok {
		delegateField = field.Value().(Delegate)
	}

	request := packet{
		Pattern:  CleanPattern(s),
		Meta:     metaField,
		Delegate: delegateField,
		Trace: Trace{
			TraceID: nuid.Next(),
		},
		Request: request{
			ID:          nuid.Next(),
			RequestType: RequestType,
		},
	}

	fmt.Println(request.Pattern)

	data, _ := jsoniter.Marshal(&request)
	m, err := h.Conn.Request(topic, data, h.Opts.Timeout*time.Millisecond)

	if err != nil {
		log.Fatal(err)
		context.Error = &Error{Message: err.Error()}
		return context
	}

	pack := packet{}
	mErr := jsoniter.Unmarshal(m.Data, &pack)

	if mErr != nil {
		log.Fatal(mErr)
		context.Error = &Error{Message: mErr.Error()}
		return context
	}

	errResMap := mapstructure.Decode(pack.Result, out)

	if errResMap != nil {
		panic(errResMap)
	}

	responseError := pack.Error

	// create container for error
	errorMsg := &Error{}

	if responseError != nil {
		// Decode error map to struct
		errErrMap := mapstructure.Decode(responseError, &errorMsg)

		if errErrMap != nil {
			panic(errErrMap)
		}
	}

	context.Trace = pack.Trace
	context.Meta = pack.Meta
	context.Delegate = pack.Delegate

	return context
}

// Dissect the cb Handler's signature
func ArgInfo(cb Handler) ([]reflect.Type, int) {
	cbType := reflect.TypeOf(cb)

	if cbType.Kind() != reflect.Func {
		panic("hemera: Handler needs to be a func")
	}

	numArgs := cbType.NumIn()
	argTypes := []reflect.Type{}

	for i := 0; i < numArgs; i++ {
		argTypes = append(argTypes, cbType.In(i))
	}

	return argTypes, numArgs
}

func CleanPattern(s *structs.Struct) interface{} {
	var pattern = make(map[string]interface{})

	for _, f := range s.Fields() {
		if f.IsExported() {
			switch f.Value().(type) {
			case Meta:
			case Delegate:
			default:
				pattern[f.Name()] = f.Value()
			}
		}

	}

	return pattern
}
