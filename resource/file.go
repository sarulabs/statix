package resource

import (
	"io/ioutil"

	"github.com/sarulabs/statix/helpers"
)

// File is a resource which content is read from a file.
type File struct {
	Path string
}

// NewFile creates a new File resource.
func NewFile(path string) *File {
	return &File{
		Path: path,
	}
}

// Dump returns the content of the resource.
func (f *File) Dump() ([]byte, error) {
	return ioutil.ReadFile(f.Path)
}

// In creates a copy of File resource with its path modified.
// If the File path is absolute, nothing changes.
// If the File path is relative, it is rewritten to be based in the path parameter.
func (f *File) In(path string) Resource {
	return &File{
		Path: helpers.RewritePath(path, f.Path),
	}
}
