package alteration

import "github.com/sarulabs/statix/resource"

// UglifyCss is an alteration that can apply uglifycss to a resource.
// Bin is the path to uglifycss executable.
type UglifyCss struct {
	Bin string
}

// NewUglifyCss creates a new UglifyCss alteration.
func NewUglifyCss(bin string) UglifyCss {
	return UglifyCss{
		Bin: bin,
	}
}

// Alter runs uglifycss on a resource returns a one.
func (ucss UglifyCss) Alter(r resource.Resource) (resource.Resource, error) {
	return ExecCommand(ucss.Bin, TmpInputFile{Resource: r})
}
