package hemera

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"
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
	ErrInvalidActHandlerArguments = errors.New("Act Handler requires at least one argument")
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
func Create(conn *nats.Conn, options ...Option) (Hemera, error) {
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
	argTypes, numArgs := argInfo(cb)

	if numArgs < 3 {
		return nil, ErrInvalidAddHandlerArguments
	}

	// Response struct
	argMsgType := argTypes[1]

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

		// Get "Value" of the reply callback for the reflection Call
		oPtr = reflect.ValueOf(oi)

		// array of arguments for the callback handler
		oV := []reflect.Value{oContextPtr, oPtr, oReplyPtr}

		cbValue.Call(oV)
	}

	return h.Conn.QueueSubscribe(topic, topic, natsCB)
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

	var pattern = make(map[string]interface{})

	// pattern contains only primitive values
	// and no meta, delegate informations
	for _, f := range s.Fields() {
		fn := f.Name()

		if !strings.HasSuffix(fn, "_") {
			fk := f.Kind()

			switch fk {
			case reflect.Struct:
			case reflect.Map:
			case reflect.Array:
			case reflect.Func:
			case reflect.Chan:
			case reflect.Slice:
			default:
				pattern[f.Name()] = f.Value()
			}
		}
	}

	argTypes, numArgs := argInfo(cb)

	if numArgs < 3 {
		return false, ErrInvalidActHandlerArguments
	}

	// Response struct
	argMsgType := argTypes[numArgs-1]

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
		argValues := []reflect.Value{oContextPtr, errVal, oPtr}
		cbValue.Call(argValues)
	} else {
		errVal := reflect.ValueOf(errorMsg)
		argValues := []reflect.Value{oContextPtr, errVal, oPtr}
		cbValue.Call(argValues)
	}

	return true, nil
}

// Dissect the cb Handler's signature
func argInfo(cb Handler) ([]reflect.Type, int) {
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
