// Package abif contains ABIF parsers and writers.
//
// ABIF (Applied Biosystems Inc., Format) is a binary key-value store format
// commonly used by Sanger sequencing machine software.
//
// This implementation is compatible with major version 1 of the ABIF format and attempts to follow
// all implementation guidelines specified in this document:
// https://projects.nfstc.org/workshops/resources/articles/ABIF_File_Format.pdf
package abif

// An ABIF is a parsed ABIF file.
type ABIF struct {
	MajorVersion int
	MinorVersion int
	Data         map[Tag]Value
}

// A Tag is a key in the ABIF key-value store.
type Tag struct {
	Name   [4]byte
	Number int
}

// A Value is a value in the ABIF key-value store.
type Value struct {
	Type        ElementType // type of the elements
	ElementSize int16       // size of a single element in Bytes
	NumElements int32       // number of elements in the value
	Bytes       []byte      // raw data, len(Bytes) = ElementSize * NumElements
}

// An ElementType represents the type of a single element in a Value.
type ElementType int16

// ElementType definitions.
const (
	Byte         = 1
	Char         = 2
	Word         = 3
	Long         = 5
	Float        = 7
	Double       = 8
	Date         = 10
	Time         = 11
	Thumb        = 12
	Bool         = 13
	PString      = 18
	CString      = 19
	RootDirEntry = 1023
)

type dirEntry struct {
	Name        [4]byte
	Number      int32
	ElementType ElementType
	ElementSize int16
	NumElements int32
	DataSize    int32
	DataOffset  int32
	DataHandle  int32
}
