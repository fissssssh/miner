package web

type RequestError struct {
	Status int
	Msg    string
}

func (e *RequestError) Error() string {
	return e.Msg
}

func NewRequestError(status int, msg string) *RequestError {
	return &RequestError{
		Status: status,
		Msg:    msg,
	}
}
