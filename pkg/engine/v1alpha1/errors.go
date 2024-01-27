package v1alpha1

type Error string

func (e Error) Error() string {
	return string(e)
}

func IsError(err error, e Error) bool {
	typed, ok := err.(Error)
	if !ok {
		return false
	}
	return typed == e
}

const (
	ErrTimeout Error = "Timeout"
)
