package resource

// Alteration is the interface defining structures that can take
// a resource and return a new modified one.
type Alteration interface {
	Alter(Resource) (Resource, error)
}

// AlteredResource is a Resource to which is applied a list of alterations.
type AlteredResource struct {
	Resource    Resource
	Alterations []Alteration
}

// NewAlteredResource creates a new AlteredResource.
func NewAlteredResource(r Resource, a ...Alteration) *AlteredResource {
	return &AlteredResource{
		Resource:    r,
		Alterations: a,
	}
}

// Dump returns the content of an AlteredResource.
func (ar *AlteredResource) Dump() ([]byte, error) {
	var err error
	r := ar.Resource
	for _, alteration := range ar.Alterations {
		r, err = alteration.Alter(r)
		if err != nil {
			return []byte{}, err
		}
	}
	return r.Dump()
}

// In returns a new AlteredResource that is the same as the structed on which
// the method is applied, but with the In method applied to its Resource attribute.
func (ar *AlteredResource) In(path string) Resource {
	return &AlteredResource{
		Resource:    ar.Resource.In(path),
		Alterations: ar.Alterations,
	}
}
