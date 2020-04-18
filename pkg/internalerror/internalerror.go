package internalerror

type InternalError struct {
	ErrorMessage string `json:"error"`
}

func (i InternalError) Error() string {
	return i.ErrorMessage
}

func New(message string) InternalError {
	return InternalError{
		ErrorMessage: message,
	}
}
