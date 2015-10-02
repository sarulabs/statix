package resource

import (
	"bytes"
	"testing"
)

func TestCreateBytes(t *testing.T) {
	b := NewBytes([]byte("test1"))

	if !bytes.Equal(b.Content, []byte("test1")) {
		t.Error("error creating Bytes Resource")
	}

	b = NewString("test2")

	if !bytes.Equal(b.Content, []byte("test2")) {
		t.Error("error creating Bytes Resource")
	}
}

func TestDumpBytes(t *testing.T) {
	b := NewBytes([]byte("test"))
	content, _ := b.Dump()

	if !bytes.Equal(content, []byte("test")) {
		t.Error("error dumping Bytes Resource")
	}
}

func TestInBytes(t *testing.T) {
	b := NewBytes([]byte("test"))
	b = b.In("whatever").(*Bytes)

	if !bytes.Equal(b.Content, []byte("test")) {
		t.Error("error while appling In to Bytes")
	}
}
