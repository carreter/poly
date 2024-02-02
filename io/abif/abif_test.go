package abif

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Roundtrip_From_File(t *testing.T) {
	file, err := os.Open("data/A_forward.ab1")
	if err != nil {
		t.Fatalf("%v", err)
	}

	firstParse, err := NewParser(file).Parse()
	if err != nil {
		t.Fatalf("%v", err)
	}

	write := firstParse.Bytes()

	secondParse, err := NewParser(bytes.NewReader(write)).Parse()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if diff := cmp.Diff(firstParse, secondParse); diff != "" {
		t.Errorf("file -> struct -> file -> struct roundtrip resulted in different structs (-first,+second): %v", diff)
	}
}

func Test_Roundtrip(t *testing.T) {
	testcases := []struct {
		name string
		data ABIF
	}{{
		name: "parses version",
		data: ABIF{
			MajorVersion: 1,
			MinorVersion: 4,
			Data:         map[Tag]Data{},
		},
	}, {
		name: "parses data contained in dataoffset",
		data: ABIF{
			MajorVersion: 1,
			MinorVersion: 4,
			Data: map[Tag]Data{
				{
					Name:   "asdf",
					Number: 4,
				}: {
					Type:        Byte,
					ElementSize: 1,
					NumElements: 1,
					Bytes:       []byte{0, 0, 0, 4},
				},
			},
		},
	}, {
		name: "parses multiple data entries all contained in dataoffset",
		data: ABIF{
			MajorVersion: 1,
			MinorVersion: 4,
			Data: map[Tag]Data{
				{
					Name:   "asdf",
					Number: 4,
				}: {
					Type:        Byte,
					ElementSize: 1,
					NumElements: 1,
					Bytes:       []byte{0, 0, 0, 4},
				},
				{
					Name:   "jkl;",
					Number: 5,
				}: {
					Type:        Byte,
					ElementSize: 1,
					NumElements: 1,
					Bytes:       []byte{0, 0, 0, 2},
				},
			},
		},
	}, {
		name: "parses data too large for dataoffset",
		data: ABIF{
			MajorVersion: 1,
			MinorVersion: 4,
			Data: map[Tag]Data{
				{
					Name:   "asdf",
					Number: 4,
				}: {
					Type:        Double,
					ElementSize: 8,
					NumElements: 1,
					Bytes:       []byte{0, 0, 0, 0, 0, 0, 0, 4},
				},
			},
		},
	}, {
		name: "parses multiple data entries all too large for dataoffset",
		data: ABIF{
			MajorVersion: 1,
			MinorVersion: 4,
			Data: map[Tag]Data{
				{
					Name:   "asdf",
					Number: 4,
				}: {
					Type:        Double,
					ElementSize: 8,
					NumElements: 1,
					Bytes:       []byte{0, 2, 8, 3, 1, 5, 4, 3},
				},
				{
					Name:   "jkl;",
					Number: 5,
				}: {
					Type:        Double,
					ElementSize: 8,
					NumElements: 1,
					Bytes:       []byte{0, 0, 2, 3, 4, 1, 0, 7},
				},
			},
		},
	}}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewParser(bytes.NewReader(tc.data.Bytes())).Parse()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.data, got); diff != "" {
				t.Errorf("struct -> bytes -> struct roundtrip resulted in different structs (-want,+got): %v", diff)
			}
		})
	}
}
