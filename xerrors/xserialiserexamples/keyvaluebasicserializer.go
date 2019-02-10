package xserialiserexamples

import (
	"bytes"
	"io"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

var basicKeyValueSeparator = []byte{' '}

type basicKeyValueSerializer struct {
	firstEntry bool
}

func (s *basicKeyValueSerializer) Keep(err error) bool {
	return !xerrors.IsFrameError(err)
}

func (s *basicKeyValueSerializer) CustomFormat(err error, b *bytes.Buffer) bool {
	if kvErr, ok := err.(keyValueError); ok {
		basicEncodeMultiKeyValue(b, kvErr.MultiKeyValue())
		return true
	}

	basicEncodeKeyValue(b, [2]string{"?", err.Error()})

	return true
}

func (s *basicKeyValueSerializer) Append(w io.Writer, b []byte) error {
	if !s.firstEntry {
		if _, err := w.Write(basicKeyValueSeparator); err != nil {
			return err
		}
	} else {
		s.firstEntry = false
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}

func (s *basicKeyValueSerializer) Reset() {
	s.firstEntry = true
}

// NewBasicKeyValueSerializer returns a serializer that prints non-frame errors in key-value form.
// If the error does not implement keyValueError it prints as ?-Error().
// Separate wrapped errors (and separate key-value pairs) are separated by a whitespace.
func NewBasicKeyValueSerializer() xerrors.Serializer {
	return &basicKeyValueSerializer{
		firstEntry: true,
	}
}
