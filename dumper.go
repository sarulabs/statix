package statix

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Dumper is the interface to dump asset content.
type Dumper interface {
	Dump(string, string, []byte) error
}

// FileDumper implements the Dumper interface to dump assets into files.
type FileDumper struct{}

// Dump write `data` in a file named `filename`
// and create a symlink named `symlink` to that file.
// If needed, directories will be created.
func (fd FileDumper) Dump(filename, symlink string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	os.Remove(symlink)

	err = os.Symlink(filename, symlink)
	if err != nil {
		return err
	}

	return nil
}
