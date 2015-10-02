package alteration

import (
	"bytes"
	"os"
	"testing"

	"github.com/sarulabs/statix/resource"
)

func TestUglifyCss(t *testing.T) {
	bin := os.Getenv("STATIX_TEST_UGLIFYCSS_BIN")

	if bin == "" {
		t.Skip("STATIX_TEST_UGLIFYCSS_BIN is not set")
	}

	s := resource.NewString(" html { color : red; } ")
	a := NewUglifyCss(bin)

	r, err := a.Alter(s)

	if err != nil {
		t.Error("could not alter resource")
	}

	content, err := r.Dump()

	if err != nil {
		t.Error("could not dump content")
	}

	expected := "html{color:red}\n"

	if !bytes.Equal(content, []byte(expected)) {
		t.Error("content dumped is not correct", string(content))
	}
}
