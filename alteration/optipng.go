package alteration

import (
	"strconv"

	"github.com/sarulabs/statix/resource"
)

// OptiPng is an alteration that can apply optipng to a resource.
// Bin is the path to optipng executable.
// Level is the optimization level (from 0 to 7)
type OptiPng struct {
	Bin   string
	Level int
}

// NewOptiPng creates a new OptiPng alteration.
func NewOptiPng(bin string, level int) OptiPng {
	return OptiPng{
		Bin:   bin,
		Level: level,
	}
}

// Alter runs optipng on a resource returns a one.
func (opng OptiPng) Alter(r resource.Resource) (resource.Resource, error) {
	return ExecCommand(opng.Bin, "-o", strconv.Itoa(opng.Level), "-out", TmpOutputFile{}, "-keep", TmpInputFile{Resource: r})
}
