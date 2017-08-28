package hemera

type Context struct {
	Meta     interface{}
	Delegate interface{}
	Trace    trace
	Error    *Error
}
