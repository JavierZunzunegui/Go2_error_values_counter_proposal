package xserialiserexamples

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

var (
	jsonKeyValueSeparator = []byte{','}
	jsonKeyValueOpen      = []byte{'{'}
	jsonKeyValueClose     = []byte{'}'}
)

type jsonKeyValueSerializer struct {
	firstEntry       bool
	remainingDepth   int
	remainingCustoms int
}

func (s *jsonKeyValueSerializer) Keep(err error) bool {
	return !xerrors.IsFrameError(err)
}

func (s *jsonKeyValueSerializer) CustomFormat(err error, b *bytes.Buffer) bool {
	if s.firstEntry {
		for auxErr := err; auxErr != nil; auxErr = xerrors.Last(xerrors.Unwrap(auxErr), s.Keep) {
			s.remainingDepth++
			if !isKeyValueError(auxErr) {
				s.remainingCustoms++
			}
		}
	}

	if kvErr, ok := err.(keyValueError); ok {
		jsonEncodeMultiKeyValue(b, kvErr.MultiKeyValue())
		return true
	}

	s.remainingCustoms--
	jsonW := json.NewEncoder(b)
	jsonEncodeKeyValue(b, jsonW, [2]string{"unknown_" + strconv.Itoa(s.remainingCustoms), err.Error()})

	return true
}

func (s *jsonKeyValueSerializer) Append(w io.Writer, b []byte) error {
	if s.firstEntry {
		if _, err := w.Write(jsonKeyValueOpen); err != nil {
			return err
		}
		s.firstEntry = false
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	s.remainingDepth--
	appending := jsonKeyValueSeparator
	if s.remainingDepth == 0 {
		appending = jsonKeyValueClose
	}

	if _, err := w.Write(appending); err != nil {
		return err
	}

	return nil
}

func (s *jsonKeyValueSerializer) Reset() {
	s.firstEntry = true
	s.remainingDepth = 0
	s.remainingCustoms = 0
}

// NewJSONKeyValueSerializer returns a serializer that prints errors in JSON.
// For errors implementing KeyValueError, it prints them as "key":"value".
// If the error does not implement KeyValueError it prints as "unknown_N":"Error()", ... "unknown_0":"Error()".
// Separate wrapped errors (and separate key-value pairs) are separated by comma.
func NewJSONKeyValueSerializer() xerrors.Serializer {
	return &jsonKeyValueSerializer{
		firstEntry:       true,
		remainingDepth:   0,
		remainingCustoms: 0,
	}
}
