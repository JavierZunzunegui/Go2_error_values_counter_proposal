// Package xserialiserexamples holds some examples of xerrors.Serializer implementations.
// These are not meant to be part of the proposal, but simply a demonstration of what Serializers may do.
//
// 3 serializers are provided:
// - frameOnlySerializer: serialises FrameErrors only, in full detail with newline and tab separators
// - basicKeyValueSerializer: serialises in a human-readable form of key-value pairs
// - jsonKeyValueSerializer: serialises in a JSON form of key-value pairs
//
// Together with the xerrors Serializers (basic colon and detail colon), the following features are demonstrated:
// - some errors are printed in some serializers but not in others
// - the serialised form of errors may be modified by the serializer
// - prefix/suffix/separators may be added by serializers
package xserialiserexamples
