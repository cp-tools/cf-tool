package cmd

import (
	cln "cf/client"
	cfg "cf/config"
	pkg "cf/packages"

	"fmt"
	"os"
	"path/filepath"
)

// RunFetch is called on running cf fetch
func (opt Opts) RunFetch() {
	// check if contest id is present
	if opt.contest == "" {
		pkg.Log.Error("No contest id found")
		return
	}
	// fetch countdown info
	pkg.Log.Info("Fetching details of " + opt.contClass + " " + opt.contest)
	dur, err := cln.FindCountdown(opt.group, opt.contest, opt.contClass, opt.link)
	pkg.PrintError(err, "Extraction of countdown failed")

	// contest not yet started
	// countdown till it starts
	if dur > 0 {
		pkg.Log.Warning("Contest hasn't started")
		pkg.Log.Info("Launching countdown to start")
		cln.StartCountdown(dur)
		// open problems page (once parsing is over)
		// page will be opened only for live rounds
		defer opt.RunOpen()
	}
	// Fetch ALL problems from contest page
	pkg.Log.Info("Fetching problems...")
	probs, err := cln.FetchProbs(opt.group, opt.contest, opt.contClass, opt.link)
	pkg.PrintError(err, "Extraction of contest problems failed")

	// Fetch all tests from problems page
	splInp, splOut, err := cln.FetchTests(opt.group, opt.contest, opt.contClass, "", opt.link)
	pkg.PrintError(err, "Failed to extract sample tests")
	// no sample tests found, try parsing from each problem
	if len(splInp) == 0 {
		pkg.Log.Warning("Failed to fetch tests from problems page")
		pkg.Log.Info("Fetching from page of every problem")
		pkg.Log.Notice("Please be patient")
		// iterate over all present problems
		for _, prob := range probs {
			// Problem isn't specified to be fetched
			if opt.problem != "" && prob != opt.problem {
				// enter blank tests (as they aren't required)
				splInp = append(splInp, make([]string, 0))
				splOut = append(splOut, make([]string, 0))
				continue
			}
			probInp, probOut, err := cln.FetchTests(opt.group,
				opt.contest, opt.contClass, prob, opt.link)
			pkg.PrintError(err, "Failed to extract sample tests of "+prob)
			// append sample tests to slice
			splInp = append(splInp, probInp...)
			splOut = append(splOut, probOut...)
			// if problem is pdf format (can't extract tests)
			if len(probInp) == 0 {
				pkg.Log.Warning("Unable to extract test(s) - " + prob)
				splInp = append(splInp, make([]string, 0))
				splOut = append(splOut, make([]string, 0))
			}
		}
	}

	// iterate over fetched problems tests
	for i, prob := range probs {
		// Problem isn't specified to be fetched
		if opt.problem != "" && prob != opt.problem {
			continue
		}
		// create problem folder
		path := opt.dirPath
		if opt.group == "" {
			path = filepath.Join(path, opt.contClass, opt.contest, prob)
		} else {
			path = filepath.Join(path, opt.contClass, opt.group, opt.contest, prob)
		}
		os.MkdirAll(path, os.ModePerm)
		// create tests
		for x := 0; x < len(splInp[i]); x++ {
			// create input file (form x.in)
			pkg.CreateFile(splInp[i][x], fmt.Sprintf("%v/%d.in", path, x))
			// create output file (form x.ans)
			pkg.CreateFile(splOut[i][x], fmt.Sprintf("%v/%d.out", path, x))
		}
		pkg.Log.Success(fmt.Sprintf("Fetched %d test(s) - %v", len(splInp[i]), prob))
		// generate code files if specified
		idx := cfg.Settings.DfltTmplt
		if cfg.Settings.GenOnFetch == true && idx != -1 {
			// create temp struct with updated problem value
			oo := opt
			oo.problem = prob
			// create template file in problem folder
			oo.GenCode(&cfg.Templates[idx], path)
		}
	}

	return
}
