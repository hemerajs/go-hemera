package hemera

type Error struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int16  `json:"code"`
}

func (e *Error) isZero() bool {
	return e.Name == "" && e.Message == "" && e.Code == 0
}
