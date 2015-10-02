package statix

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sarulabs/statix/helpers"
	"github.com/sarulabs/statix/resource"
)

// Asset is the interface for assets managed by statix Manager.
type Asset interface {
	RewritePaths(string, string) Asset
	Dump([]Filter) error
}

// AssetPack implements the Asset interface. It includes all the assets
// located in the AssetPack.Input directory. Only the files with an output (without md5 suffix)
// matching the AssetPack.Pattern are part of the AssetPack.
// When the asset is dumped with the AssetPack.Dumper, AssetPack.Filters are applied before
// writing the assets in the AssetPack.Output directory.
type AssetPack struct {
	Input       string
	Output      string
	Pattern     Pattern
	Alterations []resource.Alteration
	Dumper      Dumper
}

// RewritePaths returns a new AssetPack with updated input and output.
// More precisely, if the AssetPack.Input or AssetPack.Output is relative, it is prefixed
// by the `input` and `output` parameters.
func (ap AssetPack) RewritePaths(input, output string) Asset {
	return AssetPack{
		Input:       helpers.RewritePath(input, ap.Input),
		Output:      helpers.RewritePath(output, ap.Output),
		Pattern:     ap.Pattern,
		Alterations: ap.Alterations,
		Dumper:      FileDumper{},
	}
}

// Dump reads files contained in AssetPack.Input and dumps
// them in AssetPack.Output if they match AssetPack.Pattern.
// If some filters are defined in AssetPack.Filters, they will be
// applied before dumping the assets with the AssetPack.Dumper.
// If some filters are passed in the `filters` parameter, they will be applied just after
// filters in AssetPack.Filters.
func (ap AssetPack) Dump(filters []Filter) error {
	var r resource.Resource

	files, err := ap.InputFiles()
	if err != nil {
		return err
	}

	for _, filename := range files {
		c, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		r = resource.NewBytes(c)

		for _, a := range ap.Alterations {
			r, err = a.Alter(r)
			if err != nil {
				return err
			}
		}

		output, err := ap.OutputFile(filename, "")
		for _, f := range filters {
			if f.Pattern.Match(output) {
				r, err = f.Alteration.Alter(r)
				if err != nil {
					return err
				}
			}
		}

		c, err = r.Dump()
		if err != nil {
			return err
		}

		md5Output, err := ap.OutputFile(filename, "."+helpers.MD5(c))
		if err != nil {
			return err
		}

		err = ap.Dumper.Dump(md5Output, output, c)
		if err != nil {
			return err
		}
	}

	return nil
}

// InputFiles returns all the files contained in AssetPack.Input
// that are not directories and that match AssetPack.Pattern.
func (ap AssetPack) InputFiles() ([]string, error) {
	var walkError error
	files := []string{}

	info, err := os.Stat(ap.Input)
	if err != nil || !info.IsDir() {
		return files, fmt.Errorf("asset input `%s` is not a directory", ap.Input)
	}

	filepath.Walk(ap.Input, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			walkError = err
		}
		if info.IsDir() {
			return nil
		}
		if ap.Pattern.Match(filename) {
			files = append(files, filename)
		}
		return nil
	})

	return files, walkError
}

// OutputFile returns the absolute filename of an output based on the filename of the input.
// It also add a suffix in the basename (for example an md5 hash or a version number).
// The suffix is inserted just before the file extension and at the end of the filename
// if no extension was found.
func (ap AssetPack) OutputFile(filename string, suffix string) (string, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	input, err := filepath.Abs(ap.Input)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(filename, input) {
		return "", errors.New(ap.Input + " is not a prefix of " + filename)
	}

	out := bytes.NewBuffer(nil)
	out.WriteString(ap.Output)
	out.WriteByte(os.PathSeparator)
	out.WriteString(filename[len(input):])

	return helpers.AddFileSuffix(out.String(), suffix), nil
}

// SingleAsset implements the Asset interface.
// It includes only one asset (SingleAsset.Input) implementing the Asset
// interface located in the asset package. The asset will be dump in the SingleAsset.Output file.
type SingleAsset struct {
	Input  resource.Resource
	Output string
	Dumper Dumper
}

// RewritePaths returns a new SingleAsset with updated input and output.
// More precisely, if the SingleAsset.Input or SingleAsset.Output path is relative, it is prefixed
// by the `input` and `output` parameters.
func (sa SingleAsset) RewritePaths(input, output string) Asset {
	return SingleAsset{
		Input:  sa.Input.In(input),
		Output: helpers.RewritePath(output, sa.Output),
		Dumper: FileDumper{},
	}
}

// Dump dumps the asset defined in SingleAsset.Input.
// If some filters are passed in the `filters` parameter, they will be applied before
// dumping the asset.
func (sa SingleAsset) Dump(filters []Filter) error {
	r := sa.Input

	output, err := sa.OutputFile("")
	for _, f := range filters {
		if f.Pattern.Match(output) {
			r, err = f.Alteration.Alter(r)
			if err != nil {
				return err
			}
		}
	}

	c, err := r.Dump()
	if err != nil {
		return err
	}

	md5Output, err := sa.OutputFile("." + helpers.MD5(c))
	if err != nil {
		return err
	}

	err = sa.Dumper.Dump(md5Output, output, c)
	if err != nil {
		return err
	}

	return nil
}

// OutputFile returns the absolute filename of the SingleAsset.Output.
// It also add a suffix in the basename (for example an md5 hash or a version number).
// The suffix is inserted just before the file extension and at the end of the filename
// if no extension was found.
func (sa SingleAsset) OutputFile(suffix string) (string, error) {
	out, err := filepath.Abs(sa.Output)
	if err != nil {
		return "", err
	}
	return helpers.AddFileSuffix(out, suffix), nil
}
