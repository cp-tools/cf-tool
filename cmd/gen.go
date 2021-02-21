package cmd

import (
	cfg "cf/config"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

// RunGen is called on running cf gen
func (opt Opts) RunGen() {
	// check if any templates exist
	if len(cfg.Templates) == 0 {
		Log.Error("No configured template's exist")
		return
	}
	// index of template config to use
	idx := cfg.Settings.DfltTmplt
	if len(cfg.Templates) == 1 {
		idx = 0
	} else if idx == -1 || opt.All == true {
		// ask user to select desired template
		err := survey.AskOne(&survey.Select{
			Message: "Select template to generate:",
			Options: cfg.ListTmplts(cfg.Templates...),
		}, &idx)
		PrintError(err, "")
	}
	// create template in current folder
	currDir, _ := os.Getwd()
	opt.GenCode(&cfg.Templates[idx], currDir)
	return
}

// GenCode is to generate the code file in given path
func (opt Opts) GenCode(t *cfg.Template, path string) {
	// read template code file
	file, err := ioutil.ReadFile(t.Path)
	PrintError(err, "Failed to read template file")
	// clean template code (replace placeholders)
	e := Env{
		Contest:   opt.contest,
		Problem:   opt.problem,
		Group:     opt.group,
		ContClass: opt.contClass,
	}

	source := e.ReplPlaceholder(string(file))

	// name of file to be created
	currBase := filepath.Base(path)
	fName := fmt.Sprintf("%v${idx}%v", currBase, t.Ext)
	for idx := 0; ; idx++ {
		// idx value to replace in string
		e.Idx = strconv.Itoa(idx)
		name := e.ReplPlaceholder(fName)

		// check if file already exists
		if _, err := os.Stat(filepath.Join(path, name)); os.IsNotExist(err) {
			CreateFile(source, filepath.Join(path, name))
			Log.Notice("File " + name + " generated")
			break
		}
		Log.Warning("File " + name + " exists")
	}
	return
}
