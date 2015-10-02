package alteration

import (
	"os"
	"testing"

	"github.com/sarulabs/statix/resource"
)

func TestOptiPng(t *testing.T) {
	bin := os.Getenv("STATIX_TEST_OPTIPNG_BIN")

	if bin == "" {
		t.Skip("STATIX_TEST_OPTIPNG_BIN is not set")
	}

	input := resource.NewFile("testFiles/test.png")
	a := NewOptiPng(bin, 2)

	output, err := a.Alter(input)

	if err != nil {
		t.Error("could not alter resource")
	}

	outputBytes, err := output.Dump()

	if err != nil {
		t.Error("could not dump content")
	}

	inputBytes, _ := input.Dump()

	if len(inputBytes) < len(outputBytes) {
		t.Error("resource has not been optimized")
	}
}
