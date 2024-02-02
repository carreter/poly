package abif

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	libraryMajorVersion = 1
	signatureLength     = 4
	emptyHeaderLength   = 47 * 2
	headerLength        = 128
)

// A Parser holds data necessary for ABIF parsing.
type Parser struct {
	file io.ReadSeeker
}

// NewParser creates a Parser from an io.ReadSeeker.
func NewParser(file io.ReadSeeker) *Parser {
	return &Parser{file: file}
}

// Parse decodes an ABIF struct from the parser's underlying io.ReadSeeker.
func (p *Parser) Parse() (ABIF, error) {
	res := ABIF{}

	// Verify file signature.
	signature, err := p.readBytes(signatureLength)
	if err != nil {
		return res, err
	}
	if string(signature) != "ABIF" {
		return res, fmt.Errorf("incorrect file signature, got: %v, expected: ABIF", string(signature))
	}

	// Read in file version.
	rawVersion, err := p.readInt16()
	if err != nil {
		return res, err
	}
	res.MajorVersion = int(rawVersion / 100)
	res.MinorVersion = int(rawVersion) - res.MajorVersion*100
	if res.MajorVersion != libraryMajorVersion {
		return res, fmt.Errorf("ABIF major version %v is not supported, only major version %v is supported", res.MajorVersion, libraryMajorVersion)
	}

	// Read in root directory entry that points to the file's directory entries.
	rootDir, err := p.readDirEntry()
	if err != nil {
		return res, err
	}

	// Skip remaining Bytes in header.
	_, err = p.readBytes(emptyHeaderLength)
	if err != nil {
		return res, err
	}

	res.Data = make(map[Tag]Value)

	// Read in the directory entries.
	for i := 0; i < int(rootDir.NumElements); i++ {
		err := p.seek(int32(rootDir.ElementSize*int16(i)) + rootDir.DataOffset)
		if err != nil {
			return res, err
		}

		dir, err := p.readDirEntry()
		if err != nil {
			return res, err
		}

		tag := Tag{
			Name:   dir.Name,
			Number: int(dir.Number),
		}
		data := Value{
			Type:        dir.ElementType,
			ElementSize: dir.ElementSize,
			NumElements: dir.NumElements,
		}

		// Small data is contained in the DataOffset itself. Large data is pointed to
		// by the offset.
		if dir.DataSize <= 4 {
			data.Bytes = make([]byte, dir.DataSize)
			for i := 0; i < len(data.Bytes); i++ {
				// This magic incantation grabs the ith byte in dir.DataOffset, indexed from most to least significant.
				// For example, if dir.DataOffset = 0o0123 and i = 2, b will be 0o2.
				b := byte(dir.DataOffset >> int32(8*(4-i-1)))
				data.Bytes[i] = b
			}
		} else {
			err := p.seek(dir.DataOffset)
			if err != nil {
				return res, err
			}
			data.Bytes, err = p.readBytes(int(dir.ElementSize) * int(dir.NumElements))
			if err != nil {
				return res, err
			}
		}

		res.Data[tag] = data
	}

	return res, nil
}

func (p *Parser) readBytes(n int) ([]byte, error) {
	res := make([]byte, n)
	nRead, err := p.file.Read(res)
	if err != nil {
		return nil, err
	} else if nRead != n {
		return nil, fmt.Errorf("could only read %v of %v desired Bytes", nRead, n)
	}

	return res, nil
}

func (p *Parser) readInt16() (int16, error) {
	var res int16
	err := binary.Read(p.file, binary.BigEndian, &res)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (p *Parser) readDirEntry() (dirEntry, error) {
	var res dirEntry
	err := binary.Read(p.file, binary.BigEndian, &res)
	if err != nil {
		return dirEntry{}, err
	}

	return res, nil
}

func (p *Parser) seek(offset int32) error {
	newOffset, err := p.file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return err
	} else if newOffset != int64(offset) {
		return fmt.Errorf("could not seek to offset %v, got offset %v instead", offset, newOffset)
	}

	return nil
}
