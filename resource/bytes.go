package resource

// Bytes is a resource stored in a slice of bytes.
type Bytes struct {
	Content []byte
}

// NewString creates a Bytes resource from a string.
func NewString(content string) *Bytes {
	return &Bytes{
		Content: []byte(content),
	}
}

// NewBytes creates a Bytes resource from a slice of bytes.
func NewBytes(content []byte) *Bytes {
	return &Bytes{
		Content: content,
	}
}

// Dump returns the content of the resource.
func (s *Bytes) Dump() ([]byte, error) {
	return s.Content, nil
}

// In returns a copy of the resource.
func (s *Bytes) In(path string) Resource {
	return NewBytes(s.Content)
}
