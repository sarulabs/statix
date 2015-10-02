package resource

import (
	"bytes"
	"testing"
)

// ReverseAlteration is an alteration that can reverse the content of a resource
type ReverseAlteration struct{}

func (ra ReverseAlteration) Alter(ressource Resource) (Resource, error) {
	content, _ := ressource.Dump()

	for i := 0; i < len(content)/2; i++ {
		content[i], content[len(content)-1-i] = content[len(content)-1-i], content[i]
	}

	return NewBytes(content), nil
}

func TestDumpAlteredResource(t *testing.T) {
	ar := NewAlteredResource(
		NewString("inorder"),
		ReverseAlteration{},
	)
	content, _ := ar.Dump()

	if !bytes.Equal(content, []byte("redroni")) {
		t.Error("error dumping reversed Resource")
	}

	ar = NewAlteredResource(
		NewString("inorder"),
		ReverseAlteration{},
		ReverseAlteration{},
	)
	content, _ = ar.Dump()

	if !bytes.Equal(content, []byte("inorder")) {
		t.Error("error dumping AlteredResource with multiple alterations")
	}
}

func TestInAlteredResource(t *testing.T) {
	ar := NewAlteredResource(
		NewFile("file/path"),
		ReverseAlteration{},
	)
	ar = ar.In("base").(*AlteredResource)
	f := ar.Resource.(*File)

	if f.Path != "base/file/path" {
		t.Error("error while appling In to Collection")
	}
}
