package alteration

import "github.com/sarulabs/statix/resource"

// UglifyJs is an alteration that can apply uglifyjs to a resource.
// Bin is the path to uglifyjs executable.
type UglifyJs struct {
	Bin string
}

// NewUglifyJs creates a new UglifyJs alteration.
func NewUglifyJs(bin string) UglifyJs {
	return UglifyJs{
		Bin: bin,
	}
}

// Alter runs uglifyjs on a resource returns a one.
func (ujs UglifyJs) Alter(r resource.Resource) (resource.Resource, error) {
	return ExecCommand(ujs.Bin, TmpInputFile{Resource: r})
}
