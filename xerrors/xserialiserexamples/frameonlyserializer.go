package xserialiserexamples

import (
	"bytes"
	"io"
	"strings"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

var (
	frameOnlySeparator = []byte("\n" + strings.Repeat("\t", 100))
)

type frameOnlySerializer struct {
	tabs uint8
}

func (s *frameOnlySerializer) Keep(err error) bool {
	return xerrors.IsFrameError(err)
}

func (s *frameOnlySerializer) CustomFormat(error, *bytes.Buffer) bool {
	return false
}

func (s *frameOnlySerializer) Append(w io.Writer, b []byte) error {
	if s.tabs != 0 {
		if _, err := w.Write(frameOnlySeparator[:s.tabs+1]); err != nil {
			return err
		}
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	s.tabs++

	return nil
}

func (s *frameOnlySerializer) Reset() {
	s.tabs = 0
}

// NewFrameOnlySerializer returns a serializer that only prints frames, in full detail.
// Each frame is split into a new line and progressive tabs - the second frame is tabbed once, the second twice, etc.
func NewFrameOnlySerializer() xerrors.Serializer {
	return &frameOnlySerializer{
		tabs: 0,
	}
}
