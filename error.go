package hemera

type (
	Error struct {
		Name    string `json:"name"`
		Message string `json:"message"`
		Code    int16  `json:"code"`
	}
)

func NewError(name, message string, code int16) *Error {
	return &Error{
		Name:    name,
		Message: message,
		Code:    code,
	}
}

func NewErrorSimple(message string) *Error {
	return &Error{
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
