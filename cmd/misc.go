package cmd

import (
	cfg "cf/config"

	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
)

// Version of the current executable
const Version = "1.1.0"

// Global Variables for different UI formatting
var (
	writer = os.Stderr
	Green  = color.New(color.FgGreen)
	Blue   = color.New(color.FgBlue)
	Red    = color.New(color.FgRed)
	Yellow = color.New(color.FgYellow)

	Log struct {
		Success, Notice, Info, Error,
		Warning func(text ...interface{})
	}
	// LiveUI to print live data to terminal
	LiveUI struct {
		count int
		Start func()
		Print func(text ...string)
	}
)

func init() {
	// Initialise colored text output
	Log.Success = func(text ...interface{}) { Green.Fprintln(writer, text...) }
	Log.Notice = func(text ...interface{}) { fmt.Fprintln(writer, text...) }
	Log.Info = func(text ...interface{}) { Blue.Fprintln(writer, text...) }
	Log.Error = func(text ...interface{}) { Red.Fprintln(writer, text...) }
	Log.Warning = func(text ...interface{}) { Yellow.Fprintln(writer, text...) }

	// Initialise Live rendering output
	LiveUI.Start = func() { LiveUI.count = 0 }
	LiveUI.Print = func(text ...string) {
		// clear last count lines from terminal
		for i := 0; i < LiveUI.count; i++ {
			ansi.CursorPreviousLine(1)
			ansi.EraseInLine(2)
		}
		// count number of lines in text
		LiveUI.count = 1
		for _, str := range text {
			LiveUI.count += strings.Count(str, "\n")
			fmt.Println(str)
		}
	}
}

type (
	// Opts is struct docopt binds flag data to
	Opts struct {
		Config  bool `docopt:"config"`
		Gen     bool `docopt:"gen"`
		Open    bool `docopt:"open"`
		Fetch   bool `docopt:"fetch"`
		Test    bool `docopt:"test"`
		Submit  bool `docopt:"submit"`
		Watch   bool `docopt:"watch"`
		Pull    bool `docopt:"pull"`
		Upgrade bool `docopt:"upgrade"`

		Info []string `docopt:"<info>"`

		All    bool   `docopt:"--all"`
		File   string `docopt:"--file"`
		IgCase bool   `docopt:"--ignore-case"`
		Exp    int    `docopt:"--ignore-exp"`
		Tl     int    `docopt:"--time-limit"`
		SubCnt int    `docopt:"--submissions"`
		Handle string `docopt:"--handle"`
		Custom bool   `docopt:"--custom"`

		contest   string
		problem   string
		group     string
		contClass string
		dirPath   string
		link      url.URL
	}

	// Env are global (generic and non-genric) variables
	Env struct {
		// generic variables
		handle string `env:"${handle}"`
		date   string `env:"${date}"`
		time   string `env:"${time}"`

		// non-generic variables
		Contest   string `env:"${contest}"`
		Problem   string `env:"${problem}"`
		Group     string `env:"${group}"`
		ContClass string `env:"${contClass}"`
		Idx       string `env:"${idx}"`
		File      string `env:"${file}"`
		FileBase  string `env:"${fileBase}"`
	}
)

// FindContestData extracts contest / problem id from path
// and also determines the class (contest / gym) from the contest id
func (opt *Opts) FindContestData() {
	// path to current directory
	currPath, _ := os.Getwd()

	if len(opt.Info) == 0 {
		// no contest id given in flags. Fetch from folder path
		data := strings.Split(currPath, string(os.PathSeparator))
		data = append(data, make([]string, 10)...)
		sz := len(data) - 10

		// cleans path to return dir path to root folder
		clean := func(i int) string {
			str := filepath.Join(data[i:]...)
			return strings.TrimSuffix(currPath, str)
		}
		// find last directory matching 'Settings.WSName'
		for i := sz - 1; i >= 0; i-- {
			// current folder name matches configured WSName
			if data[i] == cfg.Settings.WSName {
				// path corresponds to contest directory
				if data[i+1] == "contest" || data[i+1] == "gym" {
					opt.contClass = data[i+1]
					opt.contest = data[i+2]
					opt.problem = data[i+3]
					currPath = clean(i)
					break
				} else if data[i+1] == "group" {
					opt.contClass = data[i+1]
					opt.group = data[i+2]
					opt.contest = data[i+3]
					opt.problem = data[i+4]
					currPath = clean(i)
					break
				}
			}
		}
	} else if _, err := url.ParseRequestURI(opt.Info[0]); err == nil {
		// url given in the flags. parse data from url
		data := strings.Split(opt.Info[0], "/")
		// prevent out-of-bounds accessing
		data = append(data, make([]string, 10)...)
		sz := len(data) - 10
		// iterate over each part of url and
		// find first part matching criteria
		for i := 0; i < sz; i++ {
			if data[i] == "contest" || data[i] == "gym" {
				opt.contClass = data[i]
				opt.contest = data[i+1]
				opt.problem = data[i+3]
				break
			} else if data[i] == "group" {
				opt.contClass = data[i]
				opt.group = data[i+1]
				opt.contest = data[i+3]
				opt.problem = data[i+5]
				break
			}
		}
	} else {
		// parse from command line args (for example, 1234 c2)
		data := append(opt.Info, make([]string, 10)...)
		if val, err := strconv.Atoi(data[0]); err == nil {
			if val <= 100000 {
				opt.contClass = "contest"
			} else {
				opt.contClass = "gym"
			}
			opt.contest = data[0]
			opt.problem = data[1]
		} else if len(data[0]) == 10 {
			// contClass is group (has length 10)
			opt.contClass = "group"
			opt.group = data[0]
			opt.contest = data[1]
			opt.problem = data[2]
		}
	}
	// convert problem id to lowercase
	opt.problem = strings.ToLower(opt.problem)
	// set path to folder containing contClass
	opt.dirPath = filepath.Join(currPath, cfg.Settings.WSName)
	// set common link to contest
	// dereference the url variable
	link, _ := url.Parse(cfg.Settings.Host)
	opt.link = *link
	if opt.contClass == "contest" || opt.contClass == "gym" {
		// not group, regular parsing
		opt.link.Path = path.Join(opt.link.Path, opt.contClass, opt.contest)
	} else if opt.contClass == "group" {
		// append group value to link
		opt.link.Path = path.Join(opt.link.Path, opt.contClass,
			opt.group, "contest", opt.contest)
	}
	return
}

// ReplPlaceholder replaces all global variables in text
// with their respective values. Non-generic are passed as map
func (e Env) ReplPlaceholder(text string) string {
	// set date/time
	e.handle = cfg.Session.Handle
	e.date = time.Now().Format("02-01-06")
	e.time = time.Now().Format("15:04:05")

	// replace string data
	repl := func(old, new string) string {
		return strings.ReplaceAll(text, old, new)
	}
	// omit ${idx} = 0
	if e.Idx == "0" {
		e.Idx = ""
	}
	// extract file name from ${file} value
	e.FileBase = strings.TrimSuffix(e.File, filepath.Ext(e.File))

	// iterate over struct and replace variables
	t := reflect.TypeOf(e)
	v := reflect.ValueOf(e)
	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("env")
		val := v.Field(i).String()
		text = repl(tag, val)
	}

	return text
}

// Prompt user to select source file to test/submit
func selSourceFile(files []string) (string, error) {
	// validate and set source file
	if len(files) == 0 {
		err := fmt.Errorf("No source files found\n" +
			"Ensure a suitable configured template exists")
		return "", err
	} else if len(files) == 1 {
		// set source file (only 1 present)
		return files[0], nil
	}
	// prompt user for code file to set as source file
	file := ""
	err := survey.AskOne(&survey.Select{
		Message: "Source file:",
		Options: files,
	}, &file)
	PrintError(err, "")
	return file, nil
}

// Prompt user to select template configuration to use
func selTmpltConfig(tmplt []cfg.Template) (*cfg.Template, error) {
	// validate and set template config
	if len(tmplt) == 0 {
		err := fmt.Errorf("No template configuration found\n" +
			"Ensure a suitable configured template exists")
		return nil, err
	} else if len(tmplt) == 1 {
		// set template configuration (only 1 present)
		return &tmplt[0], nil
	}
	// prompt user for template configuration to select
	var idx int
	err := survey.AskOne(&survey.Select{
		Message: "Template configuration:",
		Options: cfg.ListTmplts(tmplt...),
	}, &idx)
	PrintError(err, "")
	return &tmplt[idx], nil
}

// compress and return color coded verdict
func prettyVerdict(verdict string) string {
	// compress verdict to WA, TLE, MLE
	verdict = strings.ReplaceAll(verdict, "Wrong answer", "WA")
	verdict = strings.ReplaceAll(verdict, "Time limit exceeded", "TLE")
	verdict = strings.ReplaceAll(verdict, "Memory limit exceeded", "MLE")

	switch {
	case strings.HasPrefix(verdict, "TLE"):
		return Yellow.Sprint(verdict)
	case strings.HasPrefix(verdict, "MLE"):
		return Red.Sprint(verdict)
	case strings.HasPrefix(verdict, "WA"):
		return Red.Sprint(verdict)
	case strings.HasPrefix(verdict, "Pretests passed"):
		return Green.Sprint(verdict)
	case strings.HasPrefix(verdict, "Accepted"):
		return Green.Sprint(verdict)
	default:
		return verdict
	}
}

// PrintError outputs error (with custom message)
// and exits the program execution (if err != nil)
func PrintError(err error, desc string) {
	if err != nil {
		if desc != "" {
			Log.Error(desc)
		}
		Log.Error(err.Error())
		os.Exit(0)
	}
}

// CreateFile copies data to dst (create if not exists)
// Returns absolute path to destination file
func CreateFile(data, dst string) string {
	out, err := os.Create(dst)
	PrintError(err, "File "+dst+" couldn't be created!")
	defer out.Close()

	out.WriteString(data)
	return dst
}

/*
Parsing structure of problems
-----------------------------
- WSName
  - contests
    - ${contest}
      - ${problem}

  - gym
    - ${contest}
      - ${problem}

  - group
    - ${group}
      - ${contest}
        - ${problem}
*/
