package cmd

import (
	cln "cf/client"
	cfg "cf/config"
	pkg "cf/packages"

	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// RunTest is called on running `cf test`
func (opt Opts) RunTest() {
	// find code file to test
	file, err := cln.FindSourceFiles(opt.File)
	pkg.PrintError(err, "Failed to select source file")
	// find template configs to use
	t, err := cln.FindTmpltsConfig(file)
	pkg.PrintError(err, "Failed to select template configuration")

	// main testing starts here!!
	e := Env{
		Contest:   opt.contest,
		Problem:   opt.problem,
		Group:     opt.group,
		ContClass: opt.contClass,
		File:      file,
	}

	// run prescript
	if t.PreScript != "" {
		// replace placeholders in script
		script := e.ReplPlaceholder(t.PreScript)
		pkg.Log.Notice(script)
		// run script with timer of 20 secs
		_, _, err := cln.ExecScript(script, "", 1e9)
		pkg.PrintError(err, "")
	}

	if opt.Custom == false {
		// run traditional judge
		opt.tradJudge(*t, e)
	} else {
		// run interactive / special judge
		opt.spclJudge(*t, e)
	}

	// run postscript
	if t.PostScript != "" {
		// replace placeholders in script
		script := e.ReplPlaceholder(t.PostScript)
		pkg.Log.Notice(script)
		// run script with timer of 20 secs
		_, _, err := cln.ExecScript(script, "", 1e9)
		pkg.PrintError(err, "")
	}
	return
}

// tradJudge is the traditional judging process of running
// source code against input and comparing with reqd output
func (opt Opts) tradJudge(t cfg.Template, e Env) {
	// fetch test cases from current directory
	inp, out, err := cln.FindTests()
	pkg.PrintError(err, "Failed to parse sample tests")

	// run judge for each test file
	for i := 0; i < len(inp); i++ {
		// replace placeholders in script
		script := e.ReplPlaceholder(t.Script)
		// run script and calc time taken
		elapsed, stdout, err := cln.ExecScript(script, inp[i], opt.Tl)
		stdout, out[i] = cln.Validator(stdout, out[i], opt.IgCase, opt.Exp)
		// todo : add functionality to return json string of verdict
		switch {
		case elapsed.Seconds() >= float64(opt.Tl):
			// print TLE message (add support for custom time limit)
			color.Yellow("#%d: TLE .... %v\n", i, elapsed.String())

		case err != nil:
			// print RTE message with error data
			color.Red("#%d: RTE .... %v\n", i, err.Error())

		case stdout != out[i]:
			// print WA message and diff output
			color.Red("#%d: WA .... %v\n", i, elapsed.String())
			cln.PrintDiff(inp[i], stdout, out[i])

		default:
			// print AC message
			color.Green("#%d: AC .... %v\n", i, elapsed.String())
		}
	}
	return
}

func (opt Opts) spclJudge(t cfg.Template, e Env) {
	// run script in terminal
	script := e.ReplPlaceholder(t.Script)
	cmds := strings.Split(script, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	// set stdin / stdout / stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	/*
		@todo Rename this function to something more apt
		@body Judge doesn't seem to be the right word in this case
		@body as no validation of solution takes place.\n\n
		@body What would an apt term for the same be? Interactive Session?
	*/

	// inform user that interactive judge has started
	pkg.Log.Success("-----Judge begins-----\n")
	cmd.Run()
	pkg.Log.Success("\n-----Judge closed-----")

	return
}
