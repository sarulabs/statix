package statix

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sarulabs/statix/resource"
)

// ReverseAlteration is an alteration that can reverse the content of a resource
type ReverseAlteration struct{}

func (ra ReverseAlteration) Alter(ressource resource.Resource) (resource.Resource, error) {
	content, _ := ressource.Dump()

	for i := 0; i < len(content)/2; i++ {
		content[i], content[len(content)-1-i] = content[len(content)-1-i], content[i]
	}

	return resource.NewBytes(content), nil
}

func createInputFiles() {
	os.MkdirAll("./tests/in/dirIn/subDir", 0777)
	ioutil.WriteFile("./tests/in/dirIn/a1", []byte("pack-a1"), 0777)
	ioutil.WriteFile("./tests/in/dirIn/subDir/a2.ext", []byte("pack-a2"), 0777)
}

func removeTestFiles() {
	os.RemoveAll("./tests")
}

func getManagerTest() Manager {
	return Manager{
		Input:  "./tests/in",
		Output: "./tests/out",
		Server: Server{
			Directory: ".",
			URL:       "http://www.example.com/static",
		},
		Filters: []Filter{
			{
				Alteration: ReverseAlteration{},
				Pattern:    NewExtensionPattern("ext"),
			},
		},
		Assets: map[string]Asset{
			"pack": AssetPack{
				Input:  "dirIn",
				Output: "dirOut",
			},
			"single": SingleAsset{
				Output: "single.ext",
				Input:  resource.NewString("single"),
			},
		},
	}
}

func TestManagerDump(t *testing.T) {
	removeTestFiles()
	createInputFiles()
	defer removeTestFiles()

	m := getManagerTest()
	err := m.Dump()

	if err != nil {
		t.Error(err)
	}

	// testing a1 content
	content, err := ioutil.ReadFile("./tests/out/dirOut/a1")
	if err != nil {
		t.Error("could not read asset a1")
	}
	if !bytes.Equal(content, []byte("pack-a1")) {
		t.Error("a1 should contain pack-a1 instead of ", string(content))
	}

	// testing a2.ext content
	content, err = ioutil.ReadFile("./tests/out/dirOut/subDir/a2.ext")
	if err != nil {
		t.Error("could not read asset a2.ext")
	}
	if !bytes.Equal(content, []byte("2a-kcap")) {
		t.Error("a2.ext should contain 2a-kcap instead of ", string(content))
	}

	// testing single.ext content
	content, err = ioutil.ReadFile("./tests/out/single.ext")
	if err != nil {
		t.Error("could not read asset single.ext")
	}
	if !bytes.Equal(content, []byte("elgnis")) {
		t.Error("single.ext should contain elgnis instead of ", string(content))
	}
}

func TestManagerSymlink(t *testing.T) {
	removeTestFiles()
	createInputFiles()
	defer removeTestFiles()

	m := getManagerTest()
	m.Dump()

	// testing a1 symlink
	symlink, err := m.Symlink("pack", "a1")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !strings.HasSuffix(symlink, "/out/dirOut/a1") {
		t.Errorf("symlink to a1 does not have the correct filename")
	}
	filename, err := filepath.EvalSymlinks(symlink)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !strings.HasSuffix(filename, "/out/dirOut/a1.df54fa5f220b244f5ed919c871fe56f0") {
		t.Error("symlink to a1 is not correct")
	}

	// testing a2.ext symlink
	symlink, err = m.Symlink("pack", "subDir/a2.ext")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !strings.HasSuffix(symlink, "/out/dirOut/subDir/a2.ext") {
		t.Errorf("symlink to a2.ext does not have the correct filename")
	}
	filename, err = filepath.EvalSymlinks(symlink)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !strings.HasSuffix(filename, "/out/dirOut/subDir/a2.4ee925ce5ea7f2ce1d0fe37f273dff23.ext") {
		t.Error("symlink to a2.ext is not correct")
	}

	// testing single.ext symlink
	symlink, err = m.Symlink("single")
	if err != nil {
		t.Errorf(err.Error())
	}
	if !strings.HasSuffix(symlink, "/out/single.ext") {
		t.Errorf("symlink to single.ext does not have the correct filename")
	}
	filename, err = filepath.EvalSymlinks(symlink)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !strings.HasSuffix(filename, "/out/single.b34c31dcb721861cd51bfa6f3d850524.ext") {
		t.Error("symlink to single.ext is not correct")
	}
}

func TestManagerURL(t *testing.T) {
	removeTestFiles()
	createInputFiles()
	defer removeTestFiles()

	m := getManagerTest()
	m.Dump()

	var url, expected string

	// testing a1 url
	url = m.URL("pack", "a1")
	expected = "http://www.example.com/static/dirOut/a1.df54fa5f220b244f5ed919c871fe56f0"
	if url != expected {
		t.Error("url should be ", expected, " instead of ", url)
	}

	// testing a2.ext url
	url = m.URL("pack", "subDir/a2.ext")
	expected = "http://www.example.com/static/dirOut/subDir/a2.4ee925ce5ea7f2ce1d0fe37f273dff23.ext"
	if url != expected {
		t.Error("url should be ", expected, " instead of ", url)
	}

	// testing single.ext url
	url = m.URL("single")
	expected = "http://www.example.com/static/single.b34c31dcb721861cd51bfa6f3d850524.ext"
	if url != expected {
		t.Error("url should be ", expected, " instead of ", url)
	}
}
