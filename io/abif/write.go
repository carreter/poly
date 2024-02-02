package abif

import (
	"bytes"
	"encoding/binary"
)

// Bytes converts an ABIF to its binary representation.
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
	currDataOffset := headerLength + rootDir.DataSize
	for _, tag := range tags {
		data := abif.Data[tag]

		newDirEntry := dirEntry{
			Name:        tag.Name,
			Number:      int32(tag.Number),
			ElementType: data.Type,
			ElementSize: data.ElementSize,
			NumElements: data.NumElements,
			DataSize:    int32(len(data.Bytes)),
			DataHandle:  0,
		}

		// If the data itself fits in the DataOffset (an int32/4 Bytes), put it there. Otherwise,
		// write the current data offset and then increment the offset by the length of the data.
		if len(data.Bytes) <= 4 {
			// Pad with zeros on the right.
			paddedBytes := data.Bytes
			for i := 0; i < 4-len(data.Bytes); i++ {
				paddedBytes = append(paddedBytes, 0)
			}

			err := binary.Read(bytes.NewReader(paddedBytes), binary.BigEndian, &(newDirEntry.DataOffset))
			if err != nil {
				panic(err)
			}
		} else {
			newDirEntry.DataOffset = currDataOffset
			currDataOffset += int32(data.DataSize())
		}

		err := binary.Write(buf, binary.BigEndian, newDirEntry)
		if err != nil {
			panic(err)
		}
	}

	// Write the data.
	for _, tag := range tags {
		data := abif.Data[tag]
		if data.DataSize() > 4 {
			buf.Write(data.Bytes)
		}
	}

	return buf.Bytes()
}
