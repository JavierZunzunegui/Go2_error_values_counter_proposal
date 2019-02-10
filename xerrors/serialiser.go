package xerrors

import (
	"bytes"
	"io"
)

// Serializer defines how errors will be serialised.
type Serializer interface {
	// Keep returns true if the Serializer will serialise such error and false if it is being ignored.
	// It should be based on the immediate error provided, not it's wrapped errors, if any.
	Keep(error) bool

	// CustomFormat allows the Serializer to serialise the error in some way other that the error's Error() method.
	// If the error has a custom serialisation, it should be written to the buffer and return true.
	// If it doesn't, nothing should be written to the buffer and the it should return false.
	CustomFormat(error, *bytes.Buffer) bool

	// Append writes the provided bytes to the writer, along with any custom prefix and/or suffix.
	// The provided bytes will be the content of CustomFormat's buffer, or Error().
	Append(io.Writer, []byte) error

	// Reset is an implementation detail to allow Serializer to be reused to minimise memory allocations.
	// It should return the Serializer to the same estate as when it was first initialised.
	Reset()
}
