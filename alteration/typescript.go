package alteration

import "github.com/sarulabs/statix/resource"

// TypeScript is an alteration that can run the typescript compiler on a resource.
// Bin is the path to the typescript executable.
type TypeScript struct {
	Bin string
}

// NewTypeScript creates a new TypeScript.
func NewTypeScript(bin string) TypeScript {
	return TypeScript{
		Bin: bin,
	}
}

// Alter runs the typescript compiler on a resource
// and returns a compiled one.
func (ts TypeScript) Alter(r resource.Resource) (resource.Resource, error) {
	switch r := r.(type) {
	case *resource.File:
		return ExecCommand(ts.Bin, "--out", TmpOutputFile{}, r.Path)
	default:
		return ExecCommand(ts.Bin, "--out", TmpOutputFile{}, TmpInputFile{Resource: r, Suffix: ".ts"})
	}
}
