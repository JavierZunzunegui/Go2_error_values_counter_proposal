package xerrors_test

import (
	"testing"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

type panickingError struct{}

func (panickingError) Error() string { panic("I should not be called") }

func TestWrappingError_Error(t *testing.T) {
	if xerrors.Wrap("wrapper", panickingError{}).Error() != "wrapper" {
		t.Fatal("'wrappingError.Error' method must not access the wrapped error's message")
	}
}
