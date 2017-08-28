package hemera

type Context struct {
	Meta     Meta
	Delegate Delegate
	Trace    Trace
	Error    *Error
}
