package alteration

import (
	"bytes"
	"os"
	"testing"

	"github.com/sarulabs/statix/resource"
)

func TestTypeScript(t *testing.T) {
	bin := os.Getenv("STATIX_TEST_TYPESCRIPT_BIN")

	if bin == "" {
		t.Skip("STATIX_TEST_TYPESCRIPT_BIN is not set")
	}

	s := resource.NewString("class Test{}")
	a := NewTypeScript(bin)

	r, err := a.Alter(s)

	if err != nil {
		t.Error("could not alter resource", err.Error())
	}

	content, err := r.Dump()

	if err != nil {
		t.Error("could not dump content")
	}

	expected := `var Test = (function () {
    function Test() {
    }
    return Test;
})();
`

	if !bytes.Equal(content, []byte(expected)) {
		t.Error("content dumped is not correct", string(content))
	}
}

func TestTypeScriptFile(t *testing.T) {
	bin := os.Getenv("STATIX_TEST_TYPESCRIPT_BIN")

	if bin == "" {
		t.Skip("STATIX_TEST_TYPESCRIPT_BIN is not set")
	}

	s := resource.NewFile("./testFiles/test-main.ts")
	a := NewTypeScript(bin)

	r, err := a.Alter(s)

	if err != nil {
		t.Error("could not alter resource", err.Error())
	}

	content, err := r.Dump()

	if err != nil {
		t.Error("could not dump content")
	}

	expected := `var Test = (function () {
    function Test() {
    }
    return Test;
})();
/// <reference path="./test-ref.ts"/>
`

	if !bytes.Equal(content, []byte(expected)) {
		t.Error("content dumped is not correct", string(content))
	}
}
