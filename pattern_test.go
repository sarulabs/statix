package statix

import "testing"

func TestEmptyPattern(t *testing.T) {
	p := NewPattern("")

	if !p.Match("XXX") {
		t.Error("XXX should be matched")
	}

	if !p.Match("") {
		t.Error("empty string should be matched")
	}
}

func TestSingleExtension(t *testing.T) {
	p := NewExtensionPattern("js")

	if p.Match(".jsXXX") {
		t.Error(".jsXXX should not be matched")
	}

	if p.Match("XXXjs") {
		t.Error("XXXjs should not be matched")
	}

	if !p.Match(".js") {
		t.Error(".js should be matched")
	}

	if !p.Match("XXX.js") {
		t.Error("XXX.js should be matched")
	}
}

func TestMultipleExtensions(t *testing.T) {
	p := NewExtensionPattern("js", "ts")

	if p.Match("XXXjs") {
		t.Error("XXXjs should not be matched")
	}

	if p.Match("XXX.cs") {
		t.Error("XXX.cs should not be matched")
	}

	if !p.Match(".js") {
		t.Error(".js should be matched")
	}

	if !p.Match("XXX.js") {
		t.Error("XXX.js should be matched")
	}

	if !p.Match("XXX.ts") {
		t.Error("XXX.ts should be matched")
	}
}
