package cln

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// SelfUpgrade downloads latest release and overwrites current binary
// Copied from https://github.com/yitsushi/totp-cli/blob/master/command/update.go
func SelfUpgrade(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzr, _ := gzip.NewReader(resp.Body)
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	header, _ := tr.Next()

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	// directory of executable
	dir := path.Dir(exe)
	file, err := ioutil.TempFile(dir, header.Name)
	if err != nil {
		return err
	}
	defer file.Close()
	// copy data to temp file
	_, err = io.Copy(file, tr)
	if err != nil {
		return err
	}
	// set permission and replace old binary
	file.Chmod(0755)
	err = os.Rename(file.Name(), exe)
	return err
}
