package alteration

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/sarulabs/statix/helpers"
	"github.com/sarulabs/statix/resource"
)

// JpegOptim is an alteration that can apply optipng to a resource.
// Bin is the path to optipng executable.
// StripAll: strip all (Comment & Exif) markers from output file
// Max: set maximum image quality factor (from 0 to 100)
type JpegOptim struct {
	Bin      string
	StripAll bool
	Max      int
}

// NewJpegOptim creates a new JpegOptim alteration.
func NewJpegOptim(bin string, stripAll bool, max int) JpegOptim {
	return JpegOptim{
		Bin:      bin,
		StripAll: stripAll,
		Max:      max,
	}
}

// Alter runs optipng on a resource returns a one.
func (jpegOptim JpegOptim) Alter(r resource.Resource) (resource.Resource, error) {
	// create temporary file
	f, err := helpers.TempFile("", "statix_filter_", ".jpg")
	if err != nil {
		return &resource.Empty{}, err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	c, err := r.Dump()
	if err != nil {
		return &resource.Empty{}, err
	}

	_, err = f.Write(c)
	if err != nil {
		return &resource.Empty{}, err
	}

	// execute command
	var command *exec.Cmd
	bufOut := bytes.NewBuffer(nil)
	bufErr := bytes.NewBuffer(nil)
	if jpegOptim.StripAll {
		command = exec.Command(jpegOptim.Bin, "-m", strconv.Itoa(jpegOptim.Max), f.Name())
	} else {
		command = exec.Command(jpegOptim.Bin, "--strip-all", "-m", strconv.Itoa(jpegOptim.Max), f.Name())
	}
	command.Stdout = bufOut
	command.Stderr = bufErr
	err = command.Run()
	if err != nil {
		return &resource.Empty{}, fmt.Errorf("command error:\n%s", bufErr.String())
	}

	// read result
	c, err = ioutil.ReadFile(f.Name())
	if err != nil {
		return &resource.Empty{}, err
	}
	return resource.NewBytes(c), nil
}
