package xerrors

// Wrapper provides support for error wrapping - an error that contains another error.
// All errors going forward should implement Wrapper.
type Wrapper interface {
	error

	// Unwrap gives access to the internal wrapped error.
	Unwrap() error
}

// Unwrap is helper function that, if the provided error implements Wrapper, returns the internal one.
func Unwrap(err error) error {
	wErr, ok := err.(Wrapper)
	if !ok {
		return nil
	}

	return wErr.Unwrap()
}

// Last walks the input error and its wrapped ones (recursively), returning the first one that returns true on f.
// It is the main means to identify if a certain error type exists within the wrap chain of another.
// Typed helpers (that return a certain error type or more restrictive error interface) are encouraged.
func Last(err error, f func(error) bool) error {
	for ; err != nil; err = Unwrap(err) {
		if f(err) {
			return err
		}
	}

	return nil
}

func cause(err error) bool {
	return Unwrap(err) == nil
}

// Cause returns the last error in the wrap chain, defined as that which wraps no other error.
func Cause(err error) error {
	return Last(err, cause)
}

// Wrapping is a helper struct to facilitate error types to implement Wrapper.
// Embed it in an error type and it provides the Unwrap method.
type Wrapping struct {
	err error
}

// Unwrap returns the internal wrapped error
func (w Wrapping) Unwrap() error {
	return w.err
}

// NewWrapping initialises a Wrapping.
// By default it will also produce a FrameError with information about the caller of Wrap.
// This can be either disabled or the caller modified by use of optional arguments.
func NewWrapping(err error, opts ...WrapOptionFunc) Wrapping {
	return newWrapping(err, wrapOptions{skip: 1}, opts...)
}

type wrapOptions struct {
	omitFrame bool
	skip      uint8
}

// WrapOptionFunc represent optional arguments to NewWrapping or Wrap methods.
type WrapOptionFunc = func(wrapOptions) wrapOptions

// OmitFrame stops frames from being included in NewWrapping or Wrap methods.
func OmitFrame() WrapOptionFunc {
	return func(opts wrapOptions) wrapOptions {
		opts.omitFrame = true
		return opts
	}
}

// SkipNFrames can be used to have the frame reported in NewWrapping or Wrap be something other than the calling one.
func SkipNFrames(skip uint8) WrapOptionFunc {
	return func(opts wrapOptions) wrapOptions {
		opts.skip += skip
		return opts
	}
}

func newWrapping(err error, wrapOpts wrapOptions, opts ...WrapOptionFunc) Wrapping {
	for _, opt := range opts {
		wrapOpts = opt(wrapOpts)
	}

	if wrapOpts.omitFrame {
		return Wrapping{err: err}
	}

	return Wrapping{err: newFrameError(wrapOpts.skip+1, err)}
}
