package cmd

import (
	cln "cf/client"
	pkg "cf/packages"

	"fmt"
	"net/http"
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
	err = cln.SelfUpgrade(link)
	pkg.PrintError(err, "Failed to update tool")

	pkg.Log.Success("Successfully updated to v" + lVers.String())
	return
}
