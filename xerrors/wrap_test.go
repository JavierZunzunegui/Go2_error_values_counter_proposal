package xerrors_test

import (
	"reflect"
	"testing"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

// TestSkipNFrames is fragile to refactorings because it relies on line numbers.
// Do not refactor it lightly, it is placed at the top of this file to make it harder to break it accidentally.
func TestSkipNFrames(t *testing.T) {
	const newErrLine = 15
	newErr := func(optFunc ...xerrors.WrapOptionFunc) error {
		return xerrors.Wrap("whatever", nil, optFunc...)
	}

	_, _, line := xerrors.Unwrap(newErr()).(xerrors.FrameError).FrameLocation()
	if line != newErrLine {
		t.Fatalf("mismatched line, expected %d got %d", newErrLine, line)
	}

	const wrappedErrLine = 24
	_, _, line = xerrors.Unwrap(newErr(xerrors.SkipNFrames(1))).(xerrors.FrameError).FrameLocation()
	if line != wrappedErrLine {
		t.Fatalf("mismatched line, expected %d got %d", wrappedErrLine, line)
	}
}

func TestUnwrap(t *testing.T) {
	scenarios := []struct {
		name        string
		err         error
		expectedOut error
	}{
		{
			"nil", nil, nil,
		},
		{
			"nonWrapped", xerrors.New("msg"), nil,
		},
		{
			"singleWrapped",
			xerrors.Wrap("msg", xerrors.New("msg"), xerrors.OmitFrame()),
			xerrors.New("msg"),
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			out := xerrors.Unwrap(scenario.err)
			if !reflect.DeepEqual(out, scenario.expectedOut) {
				t.Fatalf("mismatched outputs, expected %q got %q", scenario.expectedOut, out)
			}
		})
	}
}

type fooError struct {
	xerrors.Wrapping
}

func (fooError) Error() string { return "foo" }

func (fooError) Foo() {}

var _ xerrors.Wrapper = fooError{}

type Fooer interface {
	Foo()
}

var _ Fooer = fooError{}

func TestLast(t *testing.T) {
	scenarios := []struct {
		name        string
		err         error
		f           func(error) bool
		expectedOut error
	}{
		{
			name:        "nil",
			err:         nil,
			f:           func(error) bool { panic("not to be called") },
			expectedOut: nil,
		},
		{
			name:        "nonWrappedLastAny",
			err:         xerrors.New("msg"),
			f:           func(error) bool { return true },
			expectedOut: xerrors.New("msg"),
		},
		{
			name:        "nonWrappedLastNone",
			err:         xerrors.New("msg"),
			f:           func(error) bool { return false },
			expectedOut: nil,
		},
		{
			name:        "wrappedLastAny",
			err:         xerrors.Wrap("wrapper", xerrors.New("msg"), xerrors.OmitFrame()),
			f:           func(error) bool { return true },
			expectedOut: xerrors.Wrap("wrapper", xerrors.New("msg"), xerrors.OmitFrame()),
		},
		{
			name:        "wrappedLastNone",
			err:         xerrors.Wrap("wrapper", xerrors.New("msg"), xerrors.OmitFrame()),
			f:           func(error) bool { return false },
			expectedOut: nil,
		},
		{
			name:        "wrappedFooLastAny",
			err:         xerrors.Wrap("wrapper", fooError{}, xerrors.OmitFrame()),
			f:           func(error) bool { return true },
			expectedOut: xerrors.Wrap("wrapper", fooError{}, xerrors.OmitFrame()),
		},
		{
			name:        "wrappedFooLastNone",
			err:         xerrors.Wrap("wrapper", fooError{}, xerrors.OmitFrame()),
			f:           func(error) bool { return false },
			expectedOut: nil,
		},
		{
			name:        "wrappedFooLastFooer",
			err:         xerrors.Wrap("wrapper", fooError{}, xerrors.OmitFrame()),
			f:           func(err error) bool { _, ok := err.(Fooer); return ok },
			expectedOut: fooError{},
		},
		{
			name:        "wrappedFooLastNotFooer",
			err:         xerrors.Wrap("wrapper", fooError{}, xerrors.OmitFrame()),
			f:           func(err error) bool { _, ok := err.(Fooer); return !ok },
			expectedOut: xerrors.Wrap("wrapper", fooError{}, xerrors.OmitFrame()),
		},
		{
			name:        "fooWrappingLastAny",
			err:         fooError{xerrors.NewWrapping(xerrors.New("msg"), xerrors.OmitFrame())},
			f:           func(error) bool { return true },
			expectedOut: fooError{xerrors.NewWrapping(xerrors.New("msg"), xerrors.OmitFrame())},
		},
		{
			name:        "fooWrappingLastNone",
			err:         fooError{xerrors.NewWrapping(xerrors.New("msg"), xerrors.OmitFrame())},
			f:           func(error) bool { return false },
			expectedOut: nil,
		},
		{
			name:        "fooWrappingLastFooer",
			err:         fooError{xerrors.NewWrapping(xerrors.New("msg"), xerrors.OmitFrame())},
			f:           func(err error) bool { _, ok := err.(Fooer); return ok },
			expectedOut: fooError{xerrors.NewWrapping(xerrors.New("msg"), xerrors.OmitFrame())},
		},
		{
			name:        "fooWrappingLastNotFooer",
			err:         fooError{xerrors.NewWrapping(xerrors.New("msg"), xerrors.OmitFrame())},
			f:           func(err error) bool { _, ok := err.(Fooer); return !ok },
			expectedOut: xerrors.New("msg"),
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			out := xerrors.Last(scenario.err, scenario.f)
			if !reflect.DeepEqual(out, scenario.expectedOut) {
				t.Fatal("mismatched outputs")
			}
		})
	}
}

func TestCause(t *testing.T) {
	scenarios := []struct {
		name        string
		err         error
		expectedOut error
	}{
		{
			name:        "nil",
			err:         nil,
			expectedOut: nil,
		},
		{
			name:        "nonWrapped",
			err:         xerrors.New("msg"),
			expectedOut: xerrors.New("msg"),
		},
		{
			name:        "nilWrapped",
			err:         xerrors.Wrap("msg", nil, xerrors.OmitFrame()),
			expectedOut: xerrors.Wrap("msg", nil, xerrors.OmitFrame()),
		},
		{
			name:        "wrapped",
			err:         xerrors.Wrap("wrapper", xerrors.New("msg"), xerrors.OmitFrame()),
			expectedOut: xerrors.New("msg"),
		},
		{
			name: "doubleWrapped",
			err: xerrors.Wrap(
				"wrapper_2",
				xerrors.Wrap("wrapper_1", xerrors.New("msg"), xerrors.OmitFrame()),
				xerrors.OmitFrame(),
			),
			expectedOut: xerrors.New("msg"),
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			out := xerrors.Cause(scenario.err)
			if !reflect.DeepEqual(out, scenario.expectedOut) {
				t.Fatal("mismatched outputs")
			}
		})
	}
}
