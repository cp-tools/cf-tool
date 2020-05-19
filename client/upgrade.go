package cln

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/blang/semver"
	"github.com/tidwall/gjson"
)

// FetchLatest determines latest version through API page of github
func FetchLatest(owner, repo string) (semver.Version, string, error) {
	// url link of github API to fetch latest version data from
	link := fmt.Sprintf("https://api.github.com/repos/%v/%v/releases/latest", owner, repo)
	resp, err := getReqBody(http.DefaultClient, link)
	if err != nil {
		return semver.Version{}, "", err
	}
	// gjson is used in pull too. So being used again here!
	latest := gjson.GetBytes(resp, "tag_name").String()
	lVers, err := semver.ParseTolerant(latest)
	releaseNotes := gjson.GetBytes(resp, "body").String()

	return lVers, releaseNotes, err
}

// SelfUpgrade downloads latest release and overwrites current binary
// Copied from https://github.com/yitsushi/totp-cli/blob/master/command/update.go
func SelfUpgrade(owner, repo, vers string) error {
	// compile link from passed parameters to fetch binary matching current build
	link := fmt.Sprintf("https://github.com/%v/%v/releases/download/v%v/cf_%v_%v.tar.gz",
		owner, repo, vers, runtime.GOOS, runtime.GOARCH)

	resp, err := http.Get(link)
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
