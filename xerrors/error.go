package xerrors

// New produces a unwrapped string error without any frame information.
// Use it to produce sentinel errors but otherwise Wrap is preferred, even with nil wrapped error.
func New(msg string) error {
	return baseError{msg}
}

type baseError struct {
	msg string
}

func (err baseError) Error() string {
	return err.msg
}

func (err baseError) Unwrap() error {
	return nil
}

func (err baseError) String() string {
	return err.msg
}

var _ Wrapper = baseError{}

// Wrap produces a simple wrapped string error.
// By default it will also produce a FrameError with information about the caller of Wrap.
// This can be either disabled or the caller modified by use of optional arguments.
func Wrap(msg string, err error, opts ...WrapOptionFunc) error {
	return &wrappingError{
		msg:      msg,
		Wrapping: newWrapping(err, wrapOptions{skip: 1}, opts...),
	}
}

type wrappingError struct {
	msg string
	Wrapping
}

func (err *wrappingError) Error() string {
	return err.msg
}

func (err *wrappingError) String() string {
	return err.msg
}

var _ Wrapper = (*wrappingError)(nil)
