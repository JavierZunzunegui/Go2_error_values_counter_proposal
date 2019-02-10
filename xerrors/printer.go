package xerrors

import (
	"bytes"
	"io"
	"sync"
)

type printerAlloc struct {
	buf bytes.Buffer
	s   Serializer
}

// NewPrinter initialises an error printer.
func NewPrinter(fFactory func() Serializer) *Printer {
	return &Printer{
		pool: sync.Pool{
			New: func() interface{} {
				return &printerAlloc{
					buf: bytes.Buffer{},
					s:   fFactory(),
				}
			},
		},
	}
}

// Printer is an error printer.
// It is thin wrapper around a Serializer factory, this fully defines the output of the printer.
type Printer struct {
	pool sync.Pool
}

// Write prints the serialised form of the given error into the given writer.
// It is safe to be called concurrently, but the error may be written to the writer in multiple calls to w.Write.
// For this reason it is advisable to it not to be called concurrently on a writer, or provide an intermediate Buffer.
func (p *Printer) Write(w io.Writer, err error) error {
	alloc := p.pool.Get().(*printerAlloc)

	if err := p.write(w, alloc.s, err, &alloc.buf); err != nil {
		// do not return to the pool
		return err
	}

	alloc.s.Reset()

	p.pool.Put(alloc)

	return nil
}

func (*Printer) write(w io.Writer, s Serializer, err error, auxiliary *bytes.Buffer) error {
	var ok bool
	var writerErr error

	for err = Last(err, s.Keep); err != nil; err = Last(Unwrap(err), s.Keep) {
		ok = s.CustomFormat(err, auxiliary)
		if !ok {
			auxiliary.WriteString(err.Error())
		}
		if writerErr = s.Append(w, auxiliary.Bytes()); writerErr != nil {
			return writerErr
		}
		auxiliary.Reset()
	}
	return nil
}
