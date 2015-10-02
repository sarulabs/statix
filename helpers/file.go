package helpers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
)

// RewritePath rewrites `currentPath`.
// If `currentPath` is absolute, it does nothing but clean it.
// But if `currentPath` is relative, it is prefixed by `basePath`.
func RewritePath(basePath, currentPath string) string {
	if basePath == "" || filepath.IsAbs(currentPath) {
		return filepath.Clean(currentPath)
	}

	b := bytes.NewBufferString(basePath)
	b.WriteByte(os.PathSeparator)
	b.WriteString(currentPath)

	return filepath.Clean(b.String())
}

// AddFileSuffix adds a suffix in a filename.
// If the filename has an extension, the suffix is added just before.
// The returned path is cleaned.
func AddFileSuffix(filename, suffix string) string {
	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	base := filepath.Base(filename)
	name := base[:len(base)-len(ext)]

	b := bytes.NewBuffer(nil)
	b.WriteString(dir)
	b.WriteByte(os.PathSeparator)
	b.WriteString(name)
	b.WriteString(suffix)
	b.WriteString(ext)

	return filepath.Clean(b.String())
}

// MD5 returns the md5 hash of an array of byte.
func MD5(content []byte) string {
	h := md5.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}
