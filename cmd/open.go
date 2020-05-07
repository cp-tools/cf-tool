package cmd

import (
	cfg "cf/config"
	pkg "cf/packages"

	"net/url"
	"os/exec"
	"path"
	"runtime"
)

// RunOpen is called on running `cf open`
func (opt Opts) RunOpen() {
	// check if contest id is present
	if opt.contest == "" {
		pkg.Log.Error("No contest id found")
		return
	}
	link, _ := url.Parse(cfg.Settings.Host)
	if opt.group == "" {
		// not group. Regular parsing
		link.Path = path.Join(link.Path, opt.contClass, opt.contest)
	} else {
		// append group value to link
		link.Path = path.Join(link.Path, "group", opt.group, "contest", opt.contest)
	}
	// open problems page (all problems)
	if opt.problem == "" {
		link.Path = path.Join(link.Path, "problems")
	} else {
		link.Path = path.Join(link.Path, "problem", opt.problem)
	}
	// open page in default browser
	browserOpen(link.String())
	return
}
func browserOpen(url string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		exec.Command("xdg-open", url).Start()
	}
	return
}
