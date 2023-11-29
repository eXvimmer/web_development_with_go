package errors

type publicError struct {
	err error
	msg string
}

func (pe publicError) Error() string {
	return pe.err.Error()
}

func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}

// Public wraps the original error with a new error that has `Public()` method
// that will return a string message that is acceptable to display to the
// public. This error can also be unwrapped using the traditional `errors`
// package approach.
func Public(err error, msg string) error {
	return publicError{err, msg}
}
