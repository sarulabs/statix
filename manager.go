package statix

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/sarulabs/statix/helpers"
	"github.com/sarulabs/statix/resource"
)

// Server maps a directory to an url.
type Server struct {
	Directory string
	URL       string
}

// Filter is the combination of an Alteration and a Pattern.
// The Pattern may allow to apply an Alteration only to some files.
type Filter struct {
	Alteration resource.Alteration
	Pattern    Pattern
}

// Manager contains the definition of your asset and can dump them.
// It can also get the url of your assets.
// - Manager.Input will be the base directory for all your assets with a
//     relative path. Assets with an absolute path will not be changed by
//     this attribute.
// - Manager.Output works the same way as Manager.Input does but is the base
//     directory in which your assets will be dumped.
// - Manager.Server defines an url for a directory. When your want to know the
//     url of a given asset, this attribute will be used.
// - Manager.Servers is used like Manager.Server. Use it if your assets are in different directories.
// - Filters contains a list of filters that will be applied to all your assets before dumping them.
// - Assets contains all your assets. The key of the map is the name of the asset.
type Manager struct {
	Input   string
	Output  string
	Server  Server
	Servers []Server
	Filters []Filter
	Assets  map[string]Asset
}

// Dump dumps all defined assets.
// If an asset path is relative, it is rewritten to be based
// in manager.Input and Manager.Output.
func (m Manager) Dump() error {
	input, err := filepath.Abs(m.Input)
	if err != nil {
		return err
	}

	output, err := filepath.Abs(m.Output)
	if err != nil {
		return err
	}

	for _, a := range m.Assets {
		err = a.RewritePaths(input, output).Dump(m.Filters)
		if err != nil {
			return err
		}
	}

	return nil
}

// URL returns the url of an asset thanks to its name
// and what is defined in Manager.Server and Manager.Servers.
// If an error occurs, an empty string is returned.
//
// If the asset is a SingleAsset, the name is enought to get its url.
// For example manager.Url("single") works.
//
// But if the asset is an AssetPack, you also need to give its path inside the output directory.
// For example, manager.Url("pack", "/js/jquery.js") will look into the output directory
// of the asset named "pack" for the file {outputDirectory}/js/jquery.js
func (m Manager) URL(assetName string, paths ...string) string {
	symlink, err := m.Symlink(assetName, paths...)
	if err != nil {
		return ""
	}
	url, _ := m.URLFromSymlink(symlink)
	return url
}

// Symlink returns the filename of an asset symlink.
// For a SingleAsset, only the name of the asset is needed.
// For AssetPack, you also need to provide the path of the file inside the output directory
func (m Manager) Symlink(assetName string, paths ...string) (string, error) {
	a, ok := m.Assets[assetName]
	if !ok {
		return "", fmt.Errorf("asset `%s` does not exist", assetName)
	}

	switch a.(type) {
	case AssetPack:
		path := ""
		if len(paths) > 0 {
			path = paths[0]
		}
		return m.AssetPackSymlink(a.(AssetPack), path)
	case SingleAsset:
		return m.SingleAssetSymlink(a.(SingleAsset))
	}

	return "", errors.New("asset should be an AssetPack or a SingleAsset")
}

// AssetPackSymlink returns the filename of an asset symlink.
// It does not check if the symlink exists.
func (m Manager) AssetPackSymlink(ap AssetPack, filename string) (string, error) {
	ap = ap.RewritePaths(m.Input, m.Output).(AssetPack)
	symlink := helpers.RewritePath(ap.Output, filename)
	return filepath.Abs(symlink)
}

// SingleAssetSymlink returns the filename of an asset symlink.
// It does not check if the symlink exists.
func (m Manager) SingleAssetSymlink(sa SingleAsset) (string, error) {
	sa = sa.RewritePaths(m.Input, m.Output).(SingleAsset)
	return sa.OutputFile("")
}

// URLFromSymlink returns the url of an asset given its symlink.
func (m Manager) URLFromSymlink(symlink string) (string, error) {
	filename, err := filepath.EvalSymlinks(symlink)
	if err != nil {
		return "", err
	}
	return m.URLFromFilename(filename)
}

// URLFromFilename returns the url of an asset given its filename.
// It checks if the server defined in Manager.Server matches the filename.
// If not it checks all the servers defined in Manager.Servers one by one.
// If a correct server is found, the filename is rewritten into an url.
// If not, an empty string is returned.
func (m Manager) URLFromFilename(filename string) (string, error) {
	filename = helpers.RewritePath(".", filename)
	filename, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	servers := append([]Server{m.Server}, m.Servers...)

	for _, s := range servers {
		dir := helpers.RewritePath(m.Output, s.Directory)
		dir, _ = filepath.Abs(dir)

		if strings.HasPrefix(filename, dir) {
			u, err := url.Parse(s.URL + "/" + filename[len(dir):])
			if err == nil {
				u.Path = filepath.Clean(u.Path)
				return u.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no server for file `%s`", filename)
}
