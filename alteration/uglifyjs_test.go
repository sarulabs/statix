package alteration

import (
	"bytes"
	"os"
	"testing"

	"github.com/sarulabs/statix/resource"
)

func TestUglifyJs(t *testing.T) {
	bin := os.Getenv("STATIX_TEST_UGLIFYJS_BIN")

	if bin == "" {
		t.Skip("STATIX_TEST_UGLIFYJS_BIN is not set")
	}

	s := resource.NewString(" console.log( \"ok\" ) ;")
	a := NewUglifyJs(bin)

	r, err := a.Alter(s)

	if err != nil {
		t.Error("could not alter resource")
	}

	content, err := r.Dump()

	if err != nil {
		t.Error("could not dump content")
	}

	expected := "console.log(\"ok\");"

	if !bytes.Equal(content, []byte(expected)) {
		t.Error("content dumped is not correct", string(content))
	}
}
