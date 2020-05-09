package main

import (
	cmd "cf/cmd"
	cfg "cf/config"
	pkg "cf/packages"

	"os"
	"path/filepath"

	"github.com/docopt/docopt-go"
)

const manPage = `
Usage:
  cf config
  cf gen    [-A]
  cf open   [<info>...] [--api]
  cf fetch  [<info>...] [--api]
  cf test   [[-i -e<e> -t<t>] | -C] [-f<f>] [--api]
  cf submit [<info>... -f<f>] [--api]
  cf watch  [<info>... -s<cnt>] [--api]
  cf pull   [<info>...] -H<handle>
  cf upgrade

Options:
  -A, --all                   force the selection menu to appear
  -f, --file <f>              specify source file to test / submit [default: *.*]
  -i, --ignore-case           omit character-case differences in output
  -e, --ignore-exp <e>        omit float differences <= 1e-<e> [default: 10]
  -t, --time-limit <t>        set time limit (secs) for each test case [default: 2] 
  -s, --submissions <cnt>     watch status of last <cnt> submissions [default: 0] 
  -H, --handle <handle>       cf handle (not email) of reqd user  
  -C, --custom                run interactive session, with input from stdin
      --api                   remove all escape sequences from output
  -h, --help                  show this screen
  -v, --version               show cli version
`

func main() {

	args, _ := docopt.ParseArgs(manPage, os.Args[1:], cmd.Version)
	// create ~/.cf/ folder
	path, _ := os.UserConfigDir()
	path = filepath.Join(path, "cf")
	os.Mkdir(path, os.ModePerm)
	// initialise default values of cf tool
	// WARNING InitSession() depends on InitSettings()
	// and all depend on InitFormat()
	cfg.InitTemplates(filepath.Join(path, "templates.json"))
	cfg.InitSettings(filepath.Join(path, "settings.json"))
	cfg.InitSession(filepath.Join(path, "sessions.json"))
	// bind data to struct holding flags
	// and extract contest type / path
	opt := cmd.Opts{}
	args.Bind(&opt)
	opt.FindContestData()
	pkg.IsAPI(opt.API)

	// run function based on subcommand
	switch {
	case opt.Config:
		opt.RunConfig()
	case opt.Gen:
		opt.RunGen()
	case opt.Open:
		opt.RunOpen()
	case opt.Fetch:
		opt.RunFetch()
	case opt.Test:
		opt.RunTest()
	case opt.Submit:
		opt.RunSubmit()
	case opt.Watch:
		opt.RunWatch()
	case opt.Pull:
		opt.RunPull()
	case opt.Upgrade:
		cmd.RunUpgrade()
	}
	return
}

/*
Global variables
  Generic:
	- ${handle}             : username of currently logged in user session
	- ${date}               : dd-mm-yy format of current date
	- ${time}               : hh:mm:ss format of current time

	Non-Generic:
	- ${contest}            : The contest id parsed from args / folder path
	- ${problem}            : The problem id parsed from args / folder path
	- ${group}              : The group id parsed from folder path / url
	- ${contClass}          : The class of the contest (contest / gym / group)
	- ${idx}                : index of iteration (eg: c${idx} as name of gen file)
	- ${file}               : file you wish to test / submit
	- ${fileBase}           : file path (without extension) you wish to test / submit
*/
