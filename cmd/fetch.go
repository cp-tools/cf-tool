package cmd

import (
	cln "cf/client"
	cfg "cf/config"

	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RunFetch is called on running cf fetch
func (opt Opts) RunFetch() {
	// check if contest id is present
	if opt.contest == "" {
		Log.Error("No contest id found")
		return
	}
	// fetch countdown info
	Log.Info("Fetching details of " + opt.contClass + " " + opt.contest)
	dur, err := cln.FindCountdown(opt.contest, opt.link)
	PrintError(err, "Extraction of countdown failed")

	// contest not yet started
	// countdown till it starts
	if dur > 0 {
		Log.Warning("Contest hasn't started")
		Log.Info("Launching countdown to start")
		startCountdown(dur)
		// open problems page (once parsing is over)
		// page will be opened only for live rounds
		defer opt.RunOpen()
	}
	// Fetch ALL problems from contest page
	Log.Info("Fetching problems...")
	probs, err := cln.FetchProbs(opt.contest, opt.link)
	PrintError(err, "Extraction of contest problems failed")

	// Fetch all tests from problems page
	splInp, splOut, err := cln.FetchTests(opt.contest, "", opt.link)
	PrintError(err, "Failed to extract sample tests")
	// no sample tests found, try parsing from each problem
	if len(splInp) == 0 {
		Log.Warning("Failed to fetch tests from problems page")
		Log.Info("Fetching from page of every problem")
		Log.Notice("Please be patient")
		// iterate over all present problems
		for _, prob := range probs {
			// Problem isn't specified to be fetched
			if opt.problem != "" && prob != opt.problem {
				// enter blank tests (as they aren't required)
				splInp = append(splInp, make([]string, 0))
				splOut = append(splOut, make([]string, 0))
				continue
			}
			probInp, probOut, err := cln.FetchTests(opt.contest, prob, opt.link)
			PrintError(err, "Failed to extract sample tests of "+prob)
			// append sample tests to slice
			splInp = append(splInp, probInp...)
			splOut = append(splOut, probOut...)
			// if problem is pdf format (can't extract tests)
			if len(probInp) == 0 {
				Log.Warning("Unable to extract test(s) - " + prob)
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
			CreateFile(splInp[i][x], fmt.Sprintf("%v/%d.in", path, x))
			// create output file (form x.ans)
			CreateFile(splOut[i][x], fmt.Sprintf("%v/%d.out", path, x))
		}
		Log.Success(fmt.Sprintf("Fetched %d test(s) - %v", len(splInp[i]), prob))
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

// startCountdown starts countdown of dur seconds
func startCountdown(dur int64) {
	// run timer till it runs out
	LiveUI.Start()
	for ; dur > 0; dur-- {
		h := fmt.Sprintf("%d:", dur/(60*60))
		m := fmt.Sprintf("0%d:", (dur/60)%60)
		s := fmt.Sprintf("0%d", dur%60)
		LiveUI.Print(h + m[len(m)-3:] + s[len(s)-2:])
		time.Sleep(time.Second)
	}
	// remove timer data from screen
	LiveUI.Print()
	return
}
