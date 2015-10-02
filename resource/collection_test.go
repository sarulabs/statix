package resource

import (
	"bytes"
	"testing"
)

func TestDumpCollection(t *testing.T) {
	s1 := NewString("string1")
	s2 := NewString("string2")
	s3 := NewString("string3")
	c1 := NewCollection(s1, s2)
	c2 := NewCollection(c1, s3)
	content, _ := c2.Dump()

	if !bytes.Equal(content, []byte("string1string2string3")) {
		t.Error("error dumping Collection Resource")
	}
}

func TestInCollection(t *testing.T) {
	f := NewFile("file/path")
	c := NewCollection(f)
	c = c.In("base").(*Collection)
	fModified := c.Resources[0].(*File)

	if fModified.Path != "base/file/path" {
		t.Error("error while appling In to Collection")
	}

	if f.Path != "file/path" {
		t.Error("In should not alter original resources")
	}
}
