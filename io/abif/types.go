package abif

type ABIF struct {
	MajorVersion int
	MinorVersion int
	Data         map[Tag]Data
}

type Tag struct {
	Name   string
	Number int
}

type Data struct {
	Type        ElementType
	ElementSize int16
	NumElements int32
	Bytes       []byte
}

type ElementType int16

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
