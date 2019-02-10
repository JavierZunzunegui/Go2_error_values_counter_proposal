package xerrors

import (
	"reflect"
)

func isNotFrameError(err error) bool {
	return !IsFrameError(err)
}

// Similar compares to errors and validates if they are logically identical.
// This involves checking all error types and Error() outputs are identical, but ignores wrapped FrameErrors.
// It is a replacement for reflect.DeepEqual(err1, err2) as the frame information will cause false negatives.
func Similar(err1, err2 error) bool {
	for err1, err2 = Last(err1, isNotFrameError), Last(err2, isNotFrameError); err1 != nil && err2 != nil; err1, err2 = Last(Unwrap(err1), isNotFrameError), Last(Unwrap(err2), isNotFrameError) {
		if t1, t2 := reflect.TypeOf(err1), reflect.TypeOf(err2); !reflect.DeepEqual(t1, t2) {
			return false
		}

		if err1.Error() != err2.Error() {
			return false
		}
	}

	if err1 != nil || err2 != nil {
		return false
	}

	return true
}

// Contains checks if err2 is logically contained within err1.
// This involves checking all wrapped error types and Error() outputs in err2 appear in err1 in identical order.
// It ignores wrapped FrameErrors altogether.
func Contains(err1, err2 error) bool {
	for err2 = Last(err2, isNotFrameError); err2 != nil; err2 = Last(Unwrap(err2), isNotFrameError) {
		t := reflect.TypeOf(err2)
		msg := err2.Error()

		err1 = Last(err1, func(err error) bool {
			return reflect.DeepEqual(reflect.TypeOf(err), t) && err.Error() == msg
		})
		if err1 == nil {
			return false
		}

		err1 = Unwrap(err1)
	}

	return true
}
