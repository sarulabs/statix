package resource

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestDumpFile(t *testing.T) {
	// create tempory file
	tmp, err := ioutil.TempFile("", "statix_filter_")
	if err != nil {
		t.Error("could not create tempory file")
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	err = ioutil.WriteFile(tmp.Name(), []byte("file"), 0777)
	if err != nil {
		t.Error("could not write in tempory file")
	}

	f := NewFile(tmp.Name())
	content, _ := f.Dump()

	if !bytes.Equal(content, []byte("file")) {
		t.Error("error dumping File Resource")
	}
}

func TestInFile(t *testing.T) {
	f := NewFile("file/path")
	fModified := f.In("base").(*File)

	if fModified.Path != "base/file/path" {
		t.Error("error while appling In to File")
	}

	if f.Path != "file/path" {
		t.Error("In should not alter original resources")
	}
}
