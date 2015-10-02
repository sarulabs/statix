package resource

// Resource is the interface defining structures that can dump their content.
// They should also have a In method to rewrite relative paths in the
// resource definition.
type Resource interface {
	Dump() ([]byte, error)
	In(string) Resource
}

// Empty is an empty resource.
type Empty struct{}

// Dump returns an empty slice of bytes.
func (e *Empty) Dump() ([]byte, error) {
	return []byte{}, nil
}

// In returns a copy of the Empty resource.
func (e *Empty) In(path string) Resource {
	return &Empty{}
}
