package cmd

import (
	pkg "cf/packages"

	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/blang/semver"
	"github.com/tidwall/gjson"
)

// RunUpgrade is called on running `cf upgrade`
func RunUpgrade() {
	// parse current version
	cVers := semver.MustParse(Version)
	// determine latest release version using github API
	link := "https://api.github.com/repos/infixint943/cf/releases/latest"
	resp, err := pkg.GetReqBody(&http.Client{}, link)
	pkg.PrintError(err, "Failed to fetch latest release")

	// check version of latest release from API resp
	latest := gjson.GetBytes(resp, "tag_name").String()
	lVers := semver.MustParse(latest[1:])
	// check if current release is same as latest release
	if cVers.GTE(lVers) {
		pkg.Log.Success(fmt.Sprintf("Current version (v%v) is the latest", cVers.String()))
		return
	}
	// new release found (fetch and print release notes)
	releaseNotes := gjson.GetBytes(resp, "body").String()
	pkg.Log.Success(fmt.Sprintf("New release (v%v) found", lVers.String()))
	pkg.Log.Notice(releaseNotes)
	fmt.Println()

	prompt := true
	err = survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Do you wish to upgrade from v%v to v%v?",
			cVers.String(), lVers.String()),
		Default: true,
	}, &prompt)
	pkg.PrintError(err, "")
	if prompt == false {
		pkg.Log.Info("Tool not upgraded")
		return
	}
	// url of tar file to download
	link = fmt.Sprintf("https://github.com/infixint943/cf/releases/download/%v/cf_%v_%v.tar.gz",
		latest, runtime.GOOS, runtime.GOARCH)

	pkg.Log.Info("Downloading update. Please wait.")
	err = selfUpdate(link)
	pkg.PrintError(err, "Failed to update tool")

	pkg.Log.Success("Successfully updated to v" + lVers.String())
	return
}

// Copied from https://github.com/yitsushi/totp-cli/blob/master/command/update.go
func selfUpdate(url string) error {
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
