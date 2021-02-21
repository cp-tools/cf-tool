package cmd

import (
	cln "cf/client"

	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// @todo Add support to pull submissions in groups
// @body Codeforces API currently returns only submissions
// @body in contests and gym. Maybe parse from submissions page (not API)

// RunPull is called on running cf pull
func (opt Opts) RunPull() {
	// fetch all submissions matching criteria
	Log.Success("Pulling submissions of: " + opt.Handle)
	// fetch all submissions matching criteria
	subs, err := cln.FetchSubs(opt.contest, opt.problem, opt.Handle)
	PrintError(err, "Failed to extract submission status")

	for _, sub := range subs {
		// fetch source code
		source, err := sub.FetchSubSource()
		if err != nil {
			Log.Error("Failed to pull source code:" + sub.Sid)
			continue
		}

		// create problem folder
		path := filepath.Join(opt.dirPath, opt.contClass, sub.Contest, sub.Problem)
		os.MkdirAll(path, os.ModePerm)

		fName := fmt.Sprintf("${problem}${idx}%v", cln.LangExt[sub.Lang])
		for idx := 0; ; idx++ {
			e := Env{
				Contest: sub.Contest,
				Problem: sub.Problem,
				Idx:     strconv.Itoa(idx),
			}
			name := e.ReplPlaceholder(fName)

			// check if file already exists
			if _, err := os.Stat(filepath.Join(path, name)); os.IsNotExist(err) {
				// relpath to cur dir
				cwd, _ := os.Getwd()
				relPath := strings.TrimPrefix(filepath.Join(path, name), cwd)

				CreateFile(source, filepath.Join(path, name))
				Log.Success(fmt.Sprintf("Fetched %v %v to .%v",
					sub.Contest, sub.Problem, relPath))
				break
			}
		}
	}
	return
}
