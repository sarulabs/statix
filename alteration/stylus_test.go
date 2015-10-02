package alteration

import (
	"bytes"
	"os"
	"testing"

	"github.com/sarulabs/statix/resource"
)

func TestStylus(t *testing.T) {
	bin := os.Getenv("STATIX_TEST_STYLUS_BIN")

	if bin == "" {
		t.Skip("STATIX_TEST_STYLUS_BIN is not set")
	}

	s := resource.NewString(`
			html
				color red
	`)
	a := NewStylus(bin)

	r, err := a.Alter(s)

	if err != nil {
		t.Error("could not alter resource", err.Error())
	}

	content, err := r.Dump()

	if err != nil {
		t.Error("could not dump content")
	}

	expected := `html {
  color: #f00;
}
`

	if !bytes.Equal(content, []byte(expected)) {
		t.Error("content dumped is not correct", string(content))
	}
}

func TestStylusFile(t *testing.T) {
	bin := os.Getenv("STATIX_TEST_STYLUS_BIN")

	if bin == "" {
		t.Skip("STATIX_TEST_STYLUS_BIN is not set")
	}

	s := resource.NewFile("./testFiles/test-main.styl")
	a := NewStylus(bin)

	r, err := a.Alter(s)

	if err != nil {
		t.Error("could not alter resource", err.Error())
	}

	content, err := r.Dump()

	if err != nil {
		t.Error("could not dump content")
	}

	expected := `html {
  color: #f00;
}
`

	if !bytes.Equal(content, []byte(expected)) {
		t.Error("content dumped is not correct", string(content))
	}
}
