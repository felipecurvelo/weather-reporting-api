package internalerror

type InternalError struct {
	Error string `json:"error"`
}

func New(message string) InternalError {
	return InternalError{
		Error: message,
	}
}
