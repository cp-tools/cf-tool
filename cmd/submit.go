package cmd

import (
	cln "cf/client"
	cfg "cf/config"
	pkg "cf/packages"

	"time"

	"github.com/gosuri/uitable"
)

// RunSubmit is called on running cf submit
func (opt Opts) RunSubmit() {
	// check if problem id is present
	if opt.problem == "" {
		pkg.Log.Error("No problem id found")
		return
	}
	// find code file to submit
	file, err := cln.FindSourceFiles(opt.File)
	pkg.PrintError(err, "Failed to select source file")
	// find template config to use
	t, err := cln.FindTmpltsConfig(file)
	pkg.PrintError(err, "Failed to select template configuration")

	// check login status
	usr, err := cln.LoggedInUsr()
	pkg.PrintError(err, "Failed to check login status")
	if usr == "" {
		// exit if no saved login configurations found
		if cfg.Session.Handle == "" || cfg.Session.Passwd == "" {
			pkg.Log.Error("No login details configured")
			pkg.Log.Notice("Configure login details through cf config")
			return
		}
		// attempt relogin
		pkg.Log.Warning("No logged in user session found")
		pkg.Log.Info("Attempting relogin: " + cfg.Session.Handle)
		status, err := cln.Relogin()
		pkg.PrintError(err, "Failed to login")
		if status == true {
			// logged in successfully
			pkg.Log.Success("Login successful")
		} else {
			pkg.Log.Error("Login failed")
			pkg.Log.Notice("Configure login details through 'cf config'")
			return
		}
	} else {
		// output handle details of current user
		// this is in else loop, since current user is already
		// being displayed during relogin above
		pkg.Log.Notice("Current user: " + usr)
	}

	// main submit code runs here
	err = cln.Submit(opt.group, opt.contest, opt.problem, t.LangID, file, opt.link)
	pkg.PrintError(err, "Failed to submit source code")
	pkg.Log.Success("Submitted")
	// watch submission verdict
	opt.watch()

	return
}

func (opt Opts) watch() {
	// infinite loop till verdicts declared
	pkg.LiveUI.Start()
	for query := opt.problem; ; {
		// query param to fetch submitted code verdict and not latest verdict in prob
		// fetch submission status from contest every second
		start := time.Now()

		data, err := cln.WatchSubmissions(opt.group, opt.contest, query, opt.link)
		pkg.PrintError(err, "Failed to extract submissions in contest.")
		sub := data[0]
		query = sub.ID

		tbl := uitable.New()
		tbl.Separator = " "
		tbl.AddRow("Verdict:", sub.Verdict)

		if sub.Waiting == "false" {
			tbl.AddRow("Memory:", sub.Memory)
			tbl.AddRow("Time:", sub.Time)
			pkg.LiveUI.Print(tbl.String())
			break
		}
		pkg.LiveUI.Print(tbl.String())
		// sleep for 1 second
		time.Sleep(time.Second - time.Since(start))
	}
	return
}
