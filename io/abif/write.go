package abif

import (
	"bytes"
	"encoding/binary"
)

func (abif ABIF) Bytes() []byte {
	buf := &bytes.Buffer{}

	// Write the header.
	buf.Write([]byte("ABIF")) // File signature.

	versionInt := int16(abif.MajorVersion*100 + abif.MinorVersion) // Version number.
	err := binary.Write(buf, binary.BigEndian, versionInt)
	if err != nil {
		panic(err)
	}

	rootDir := dirEntry{ // Root directory entry that points to other directory entries.
		Name:        [4]byte{'t', 'd', 'i', 'r'},
		Number:      1,
		ElementType: RootDirEntry,
		ElementSize: int16(binary.Size(dirEntry{})),
		NumElements: int32(len(abif.Data)),
		DataSize:    0,
		DataOffset:  headerLength,
		DataHandle:  0,
	}
	rootDir.DataSize = int32(rootDir.ElementSize) * rootDir.NumElements
	err = binary.Write(buf, binary.BigEndian, rootDir)
	if err != nil {
		panic(err)
	}

	buf.Write(bytes.Repeat([]byte{0}, emptyHeaderLength))

	// Generate ordered list of tags to iterate through. Necessary to ensure
	// directory entries and data are written in the same order. The ordering
	// of the list itself does not matter, just that it is consistent.
	tags := make([]Tag, 0)
	for tag := range abif.Data {
		tags = append(tags, tag)
	}

	// Write the directory entries.
	dataOffset := headerLength + rootDir.DataSize
	for _, tag := range tags {
		data := abif.Data[tag]
		offsetVal := dataOffset

		// If the data itself fits in the DataOffset (an int32/4 bytes), put it there instead of the offset value.
		if len(data.Bytes) <= 4 {
			tmp := binary.BigEndian.Uint32(data.Bytes)
			offsetVal = int32(tmp)
		} else {
			dataOffset += int32(len(data.Bytes))
		}

		newDirEntry := dirEntry{
			Name:        [4]byte{tag.Name[0], tag.Name[1], tag.Name[2], tag.Name[3]},
			Number:      int32(tag.Number),
			ElementType: data.Type,
			ElementSize: data.ElementSize,
			NumElements: data.NumElements,
			DataSize:    int32(len(data.Bytes)),
			DataOffset:  offsetVal,
			DataHandle:  0,
		}
		err := binary.Write(buf, binary.BigEndian, newDirEntry)
		if err != nil {
			panic(err)
		}
	}

	// Write the data.
	for _, tag := range tags {
		data := abif.Data[tag]
		if int(data.NumElements)*int(data.ElementSize) > 4 {
			buf.Write(data.Bytes)
		}
	}

	return buf.Bytes()
}
