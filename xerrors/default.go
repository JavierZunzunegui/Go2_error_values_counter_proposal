package xerrors

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"sync"
)

var (
	colonSeparator = []byte(": ")
	frameOpen      = []byte("(")
	frameClose     = []byte(")")
)

type colonSerializer struct {
	firstEntry bool
	keepFrames bool
	isFrame    bool
}

func (s *colonSerializer) Keep(err error) bool {
	return s.keepFrames || !IsFrameError(err)
}

func (s *colonSerializer) CustomFormat(err error, buf *bytes.Buffer) bool {
	if !s.keepFrames {
		return false
	}

	frameErr, ok := err.(FrameError)
	if !ok {
		return false
	}

	s.isFrame = true

	function, file, line := frameErr.FrameLocation()
	if i := strings.LastIndexByte(function, '/'); i != -1 {
		function = function[i+1:]
	}
	if i := strings.LastIndexByte(file, '/'); i != -1 {
		file = file[i+1:]
	}

	buf.WriteString(function)
	buf.WriteString(":")
	buf.WriteString(file)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(line))

	return true
}

func (s *colonSerializer) Append(w io.Writer, msg []byte) error {
	var err error

	if s.firstEntry {
		s.firstEntry = false
	} else {
		if s.isFrame {
			if _, err = w.Write(frameOpen); err != nil {
				return err
			}
		} else {
			if _, err = w.Write(colonSeparator); err != nil {
				return err
			}
		}
	}

	_, err = w.Write(msg)

	if s.isFrame {
		if _, err = w.Write(frameClose); err != nil {
			return err
		}
		s.isFrame = false
	}

	return err
}

func (s *colonSerializer) Reset() {
	s.firstEntry = true
	s.isFrame = false
}

func newColonSerializer(keepFrames bool) Serializer {
	return &colonSerializer{
		firstEntry: true,
		keepFrames: keepFrames,
		isFrame:    false,
	}
}

// NewColonBasicSerializer provides a formatter that appends messages with ': ' and omits frames.
// It is the serializer used by the %s representation of errors.
func NewColonBasicSerializer() Serializer {
	return newColonSerializer(false)
}

// NewColonDetailedSerializer provides a formatter that appends messages with ': '.
// Frames are printed in a shortened mode between brackets.
// It is the serializer used by the %v representation of errors.
func NewColonDetailedSerializer() Serializer {
	return newColonSerializer(true)
}

var (
	defaultPrinter         = NewPrinter(NewColonBasicSerializer)
	defaultDetailedPrinter = NewPrinter(NewColonDetailedSerializer)

	defaultEncodeBufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

// Bytes serialises to bytes an error using the default implementation of type NewColonBasicSerializer.
func Bytes(err error) []byte {
	return encodeBytes(err, defaultPrinter)
}

// DetailBytes serialises to bytes an error using the default detailed implementation of type NewColonDetailedSerializer.
func DetailBytes(err error) []byte {
	return encodeBytes(err, defaultDetailedPrinter)
}

func encodeBytes(err error, p *Printer) []byte {
	buf := defaultEncodeBufferPool.Get().(*bytes.Buffer)

	// never errors
	_ = p.Write(buf, err)

	contents := buf.Bytes()
	out := make([]byte, len(contents))
	copy(out, contents)

	buf.Reset()
	defaultEncodeBufferPool.Put(buf)

	return out
}

// String serialises an error using the default implementation of type NewColonBasicSerializer.
func String(err error) string {
	return encodeString(err, defaultPrinter)
}

// DetailString serialises an error using the default detailed implementation of type NewColonDetailedSerializer.
func DetailString(err error) string {
	return encodeString(err, defaultDetailedPrinter)
}

func encodeString(err error, p *Printer) string {
	buf := defaultEncodeBufferPool.Get().(*bytes.Buffer)

	// never errors
	_ = p.Write(buf, err)

	out := buf.String()

	buf.Reset()
	defaultEncodeBufferPool.Put(buf)

	return out
}
