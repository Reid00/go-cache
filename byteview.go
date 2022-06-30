package gocache

// A ByteView holds an immutable view of bytes.
type ByteView struct {
	b []byte
}

// Len return the view's length
func (b ByteView) Len() int {
	return len(b.b)
}

// ByteSlice return a deep copy of the data as byte slice
func (b ByteView) ByteSlice() []byte {
	return deepCopyBytes(b.b)
}

// String return the data as string
func (b ByteView) String() string {
	return string(b.b)
}

func deepCopyBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
