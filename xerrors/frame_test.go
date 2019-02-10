package xerrors_test

import (
	"strings"
	"testing"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

func TestFrameError_FrameLocation(t *testing.T) {
	// this is very fragile, it's the nature of frames. It should be the line number immediately below
	const expectedLine = 13
	err := xerrors.NewWrapping(nil).Unwrap()

	function, file, line := err.(xerrors.FrameError).FrameLocation()

	const expectedFunction = "github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors_test.TestFrameError_FrameLocation"
	if function != expectedFunction {
		t.Fatalf("mismatched function name output, expected %q got %q", expectedFunction, function)
	}

	const expectedFileSuffix = "github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors/frame_test.go"
	if !strings.HasSuffix(file, expectedFileSuffix) {
		t.Fatalf("mismatched file name output, expected suffix %q got %q", expectedFileSuffix, function)
	}

	if line != expectedLine {
		t.Fatalf("mismatched line number output, expected %d got %d", expectedLine, line)
	}
}
