package cmd

import (
	cln "cf/client"
	cfg "cf/config"
	pkg "cf/packages"

	"errors"
	"net/url"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/go-homedir"
)

// RunConfig is called on running cf config
func (opt Opts) RunConfig() {
	var choice int

	err := survey.AskOne(&survey.Select{
		Message: "Select configuration:",
		Options: []string{
			"Login to codeforces",
			"Add new code template",
			"Remove code template",
			"Other misc preferences",
		},
	}, &choice, survey.WithValidator(survey.Required))
	pkg.PrintError(err, "")

	switch choice {
	case 0:
		login()
	case 1:
		addTmplt()
	case 2:
		remTmplt()
	case 3:
		miscPrefs()
	}
	return
}

func login() {
	// check if logged in user exists
	if cfg.Session.Handle != "" {
		pkg.Log.Success("Current user: " + cfg.Session.Handle)
		pkg.Log.Warning("Current session will be overwritten")
	}
	// take input of username / password
	creds := struct{ Usr, Passwd string }{}
	err := survey.Ask([]*survey.Question{
		{
			Name:     "usr",
			Prompt:   &survey.Input{Message: "Username:"},
			Validate: survey.Required,
		}, {
			Name:     "passwd",
			Prompt:   &survey.Password{Message: "Password:"},
			Validate: survey.Required,
		},
	}, &creds)
	pkg.PrintError(err, "")
	// login and check login status
	pkg.Log.Info("Logging in")
	flag, err := cln.Login(creds.Usr, creds.Passwd)
	pkg.PrintError(err, "Login failed")
	// login was successful
	if flag == true {
		pkg.Log.Success("Login successful")
		pkg.Log.Notice("Welcome " + cfg.Session.Handle)
	} else {
		// login failed
		pkg.Log.Error("Login failed",
			"Check credentials and retry")
	}
	return
}

func addTmplt() {
	var lName []string
	for name := range cln.LangID {
		lName = append(lName, name)
	}
	pkg.Log.Info("For detailed instructions, read https://github.com/infixint943/cf/wiki/Configuration")
	tmplt := cfg.Template{}
	err := survey.Ask([]*survey.Question{
		{
			Name: "langname",
			Prompt: &survey.Select{
				Message: "Template language:",
				Options: lName,
			},
			Validate: survey.Required,
		}, {
			Name: "path",
			Prompt: &survey.Input{
				Message: "Path to code template:",
				Help: "The (relative/absolute) path to the template file you wish to use.\n" +
					"Example of valid paths are : ~/Documents/default.cpp on linux/macOS\n" +
					"and C:\\Users\\username\\Documents\\tmplt.py on windows",
			},
			Validate: func(ans interface{}) error {
				path, _ := homedir.Expand(ans.(string))
				file, err := os.Stat(path)
				if err != nil || file.IsDir() == true {
					return errors.New("path doesn't correspond to valid file")
				}
				return nil
			},
			Transform: func(ans interface{}) (newAns interface{}) {
				path, _ := homedir.Expand(ans.(string))
				return path
			},
		}, {
			Name: "alias",
			Prompt: &survey.Input{
				Message: "Template Alias:",
				Help: "A (unique) name by which you wish to recognize this template\n" +
					"For example, 'Default (C++)', 'FFT Template', etc\n",
			},
			Validate: func(ans interface{}) error {
				isPres := false
				for _, t := range cfg.Templates {
					if t.Alias == ans.(string) {
						isPres = true
						break
					}
				}
				if ans.(string) == "" {
					return errors.New("Value is required")
				} else if isPres == true {
					return errors.New("Template with same alias exists")
				}
				return nil
			},
		}, {
			Name: "prescript",
			Prompt: &survey.Input{
				Message: "Pre-script:",
				Help: "Script to run (once) to compile source file\n" +
					"For example, 'g++ -Wall ${file}', 'javac ${file}', etc\n" +
					"For details on placeholders, visit wiki documentation\n" +
					"Can be left blank, if source file can be run without compiling",
			},
		}, {
			Name: "script",
			Prompt: &survey.Input{
				Message: "Script:",
				Help: "Script to run binary/source file against test cases\n" +
					"For example, './a.out', 'java ${fileBasename}', 'python ${file}' etc\n" +
					"Field is required; Will be run once for each sample test case",
			},
			Validate: survey.Required,
		}, {
			Name: "postscript",
			Prompt: &survey.Input{
				Message: "Post-script:",
				Help: "Script to cleanup any residual binary files, etc\n" +
					"Run (once) after testing of all sample tests has finished\n" +
					"For example, 'rm a.out', 'del ${fileBasename}', etc\n" +
					"Can be left blank, if cleanup is required/desired",
			},
		},
	}, &tmplt)
	pkg.PrintError(err, "")
	// set ext and langid values manually
	tmplt.Ext = filepath.Ext(tmplt.Path)
	tmplt.LangID = cln.LangID[tmplt.LangName]
	// append new template data and save it
	cfg.Templates = append(cfg.Templates, tmplt)
	cfg.SaveTemplates()

	pkg.Log.Success("Template saved successfully")
	return
}

func remTmplt() {
	// check if any templates are present
	sz := len(cfg.Templates)
	if sz == 0 {
		pkg.Log.Error("No configured template's exist")
		return
	}

	var idx int
	err := survey.AskOne(&survey.Select{
		Message: "Template you want to remove:",
		Options: cfg.ListTmplts(-1),
	}, &idx)
	pkg.PrintError(err, "")
	// delete the template from the slice
	// and reconfigure default template settings
	cfg.Templates = append(cfg.Templates[:idx], cfg.Templates[idx+1:]...)
	if cfg.Settings.DfltTmplt == idx {
		pkg.Log.Warning("Default template configurations reset")
		cfg.Settings.DfltTmplt = -1
		cfg.Settings.GenOnFetch = false
		cfg.SaveSettings()
	}
	cfg.SaveTemplates()

	pkg.Log.Success("Templated removed successfully")
	return
}

func miscPrefs() {
	var choice int
	err := survey.AskOne(&survey.Select{
		Message: "Select configuration:",
		Options: []string{
			"Set default template",
			"Run gen after fetch",
			"Set host domain",
			"Set proxy",
		},
	}, &choice)
	pkg.PrintError(err, "")

	switch choice {
	case 0:
		// set default template
		err := survey.AskOne(&survey.Select{
			Message: "Select template",
			Options: append([]string{"None"}, cfg.ListTmplts(-1)...),
		}, &cfg.Settings.DfltTmplt)
		cfg.Settings.DfltTmplt--
		pkg.PrintError(err, "")
	case 1:
		// set GenOnFetch
		err := survey.AskOne(&survey.Confirm{
			Message: "Run gen after fetch?",
			Help: "If set to true, default template will be created for each fetched problem.\n" +
				"Default template has to be configured for this feature to work",
			Default: false,
		}, &cfg.Settings.GenOnFetch)
		pkg.PrintError(err, "")

	case 2:
		// set host domain
		err := survey.AskOne(&survey.Input{
			Message: "Url of host:",
			Help: "Host codeforces domain to fetch data from\n" +
				"Current host: " + cfg.Settings.Host,
			Default: "https://codeforces.com",
		}, &cfg.Settings.Host, survey.WithValidator(func(ans interface{}) error {
			_, err := url.ParseRequestURI(ans.(string))
			return err
		}))
		pkg.PrintError(err, "")

	case 3:
		// validate and set proxy
		err := survey.AskOne(&survey.Input{
			Message: "Proxy url:",
			Help: "Set a new proxy (should match protocol://host[:port])\n" +
				"Leave blank to reset to environment proxy\n" +
				"Current proxy: " + cfg.Settings.Proxy,
			Default: "",
		}, &cfg.Settings.Proxy, survey.WithValidator(func(ans interface{}) error {
			// reset to environment proxy
			if ans.(string) == "" {
				return nil
			}
			_, err := url.ParseRequestURI(ans.(string))
			return err
		}))
		pkg.PrintError(err, "")
	}
	cfg.SaveSettings()

	pkg.Log.Success("Configurations successfully set")
	return
}
