package cmd

import (
	cln "cf/client"

	"github.com/AlecAivazis/survey/v2"
	"github.com/blang/semver"
)

// RunUpgrade is called on running `cf upgrade`
func RunUpgrade() {
	// parse current version from set version number
	cVers, _ := semver.ParseTolerant(Version)
	lVers, releaseNotes, err := cln.FetchLatest("cp-tools", "cf")
	PrintError(err, "Failed to fetch latest release")
	// check if current release is same as latest release
	if cVers.GTE(lVers) {
		Log.Success("Current version (v" + cVers.String() + ") is the latest")
		return
	}
	// new release found (fetch and print release notes)
	Log.Success("New release (v" + lVers.String() + ") found")
	Log.Notice(releaseNotes, "\n")

	prompt := true
	err = survey.AskOne(&survey.Confirm{
		Message: "Upgrade from v" + cVers.String() + " to v" + lVers.String(),
		Default: true,
	}, &prompt)
	PrintError(err, "")

	if prompt == false {
		Log.Info("Tool not upgraded")
		return
	}

	Log.Info("Downloading update. Please wait.")
	err = cln.SelfUpgrade("cp-tools", "cf", lVers.String())
	PrintError(err, "Failed to update tool")

	Log.Success("Successfully updated to v" + lVers.String())
	return
}
