package xserialiserexamples

import (
	"bytes"
	"encoding/json"

	"github.com/JavierZunzunegui/Go2_error_values_counter_proposal/xerrors"
)

type keyValueError interface {
	xerrors.Wrapper
	MultiKeyValue() [][2]string
}

func isKeyValueError(err error) bool {
	_, ok := err.(keyValueError)
	return ok
}

type basicKeyValueError struct {
	multiKeyValue [][2]string
	xerrors.Wrapping
}

func (err *basicKeyValueError) Error() string {
	buf := bytes.Buffer{}
	basicEncodeMultiKeyValue(&buf, err.multiKeyValue)
	return buf.String()
}

func (err *basicKeyValueError) MultiKeyValue() [][2]string {
	return err.multiKeyValue
}

// NewBasicKeyValueError returns a keyValueError with default BasicKeyValueSerializer-like Error().
func NewBasicKeyValueError(multiKeyValue [][2]string, err error) error {
	return &basicKeyValueError{
		multiKeyValue: multiKeyValue,
		Wrapping:      xerrors.NewWrapping(err, xerrors.OmitFrame()),
	}
}

func basicEncodeKeyValue(buf *bytes.Buffer, kv [2]string) {
	buf.WriteString(kv[0])
	buf.WriteString("-")
	buf.WriteString(kv[1])
}

func basicEncodeMultiKeyValue(buf *bytes.Buffer, kvs [][2]string) {
	basicEncodeKeyValue(buf, kvs[0])

	for i := 1; i < len(kvs); i++ {
		buf.WriteString(" ")
		basicEncodeKeyValue(buf, kvs[i])
	}
}

func jsonEncodeKeyValue(buf *bytes.Buffer, jsonW *json.Encoder, kv [2]string) {
	jsonW.Encode(kv[0])
	buf.Truncate(buf.Len() - 1) // Encode adds \n
	buf.WriteString(":")
	jsonW.Encode(kv[1])
	buf.Truncate(buf.Len() - 1) // Encode adds \n
}

func jsonEncodeMultiKeyValue(buf *bytes.Buffer, kvs [][2]string) {
	jsonW := json.NewEncoder(buf)
	for i, kv := range kvs {
		jsonEncodeKeyValue(buf, jsonW, kv)
		if i != len(kvs)-1 {
			buf.WriteString(",")
		}
	}
}

type jsonKeyValueError struct {
	multiKeyValue [][2]string
	xerrors.Wrapping
}

func (err *jsonKeyValueError) Error() string {
	buf := bytes.Buffer{}
	buf.WriteString("{")
	jsonEncodeMultiKeyValue(&buf, err.multiKeyValue)
	buf.WriteString("}")
	return buf.String()
}

func (err *jsonKeyValueError) MultiKeyValue() [][2]string {
	return err.multiKeyValue
}

// NewJSONKeyValueError returns a KeyValueError with default JSONKeyValueSerializer-like Error().
func NewJSONKeyValueError(multiKeyValue [][2]string, err error) error {
	return &jsonKeyValueError{
		multiKeyValue: multiKeyValue,
		Wrapping:      xerrors.NewWrapping(err, xerrors.OmitFrame()),
	}
}
