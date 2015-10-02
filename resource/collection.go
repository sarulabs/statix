package resource

import "bytes"

// Collection is a list of resource.
type Collection struct {
	Resources []Resource
}

// NewCollection creates a new Collection.
func NewCollection(resources ...Resource) *Collection {
	return &Collection{
		Resources: resources,
	}
}

// Dump returns the concatenation of the content of the resources
// belonging to the collection.
func (c *Collection) Dump() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	for _, r := range c.Resources {
		c, err := r.Dump()
		if err != nil {
			return []byte{}, err
		}
		buf.Write(c)
	}
	return buf.Bytes(), nil
}

// In returns a copy of the Collection with the In method applied to
// all resources in the collection.
func (c *Collection) In(path string) Resource {
	clone := &Collection{
		Resources: []Resource{},
	}
	for _, r := range c.Resources {
		clone.Resources = append(clone.Resources, r.In(path))
	}
	return clone
}
