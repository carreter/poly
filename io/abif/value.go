package abif

// DataSize returns the size of the data in bytes.
// Notably, this is the size of the data described by the Value,
// not the size of the Value itself.
func (v Value) DataSize() uint {
	return uint(v.ElementSize) * uint(v.NumElements)
}
