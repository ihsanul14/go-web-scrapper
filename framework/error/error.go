package error

type Error struct {
	Code    int
	Message error
}

func NewError(code int, message error) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
