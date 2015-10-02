package alteration

import "github.com/sarulabs/statix/resource"

// Stylus is an alteration that can run the stylus compiler on a resource.
// Bin is the path to the stylus executable.
type Stylus struct {
	Bin string
}

// NewStylus creates a new Stylus.
func NewStylus(bin string) Stylus {
	return Stylus{
		Bin: bin,
	}
}

// Alter runs the stylus compiler on a resource and returns a compiled one.
func (ts Stylus) Alter(r resource.Resource) (resource.Resource, error) {
	switch r := r.(type) {
	case *resource.File:
		return ExecCommand(ts.Bin, "-o", TmpOutputFile{Suffix: ".css"}, r.Path)
	default:
		return ExecCommand(ts.Bin, "-o", TmpOutputFile{Suffix: ".css"}, TmpInputFile{Resource: r, Suffix: ".styl"})
	}
}
