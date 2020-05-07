package cmd

import (
	cfg "cf/config"

	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Version of the current executable
const Version = "0.9.0"

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
	}
)

// @todo Rename this file (to something more sensible)

// FindContestData extracts contest / problem id from path
// and also determines the class (contest / gym) from the contest id
func (opt *Opts) FindContestData() {
	// path to current directory
	path, _ := os.Getwd()

	if len(opt.Info) == 0 {
		// no contest id given in flags. Fetch from folder path
		data := strings.Split(path, string(os.PathSeparator))
		data = append(data, make([]string, 10)...)
		sz := len(data) - 10

		// cleans path to return dir path to root folder
		clean := func(i int) string {
			str := filepath.Join(data[i:]...)
			return strings.TrimSuffix(path, str)
		}
		// find last directory matching (contests/gym/group)
		for i := sz - 1; i >= 0; i-- {
			// path corresponds to contest directory
			if data[i] == "contest" {
				opt.contClass = "contest"
				opt.contest = data[i+1]
				opt.problem = data[i+2]
				path = clean(i)
				break
			} else if data[i] == "gym" {
				opt.contClass = "gym"
				opt.contest = data[i+1]
				opt.problem = data[i+2]
				path = clean(i)
				break
			} else if data[i] == "group" {
				opt.contClass = "group"
				opt.group = data[i+1]
				opt.contest = data[i+2]
				opt.problem = data[i+3]
				path = clean(i)
				break
			}
		}
	} else if _, err := url.ParseRequestURI(opt.Info[0]); err == nil {
		// url given in the flags. parse data from url
		data := strings.Split(opt.Info[0], "/")
		// prevent out-of-bounds accessing
		data = append(data, make([]string, 10)...)
		sz := len(data) - 10
		// iterate over each part of url
		for i := 0; i < sz; i++ {
			if data[i] == "contest" {
				opt.contClass = "contest"
				opt.contest = data[i+1]
				opt.problem = data[i+3]
				break
			} else if data[i] == "gym" {
				opt.contClass = "gym"
				opt.contest = data[i+1]
				opt.problem = data[i+3]
				break
			} else if data[i] == "group" {
				opt.contClass = "group"
				opt.group = data[i+1]
				opt.contest = data[i+3]
				opt.problem = data[i+5]
				break
			}
		}
	} else {
		// parse from command line args (like v0.2.2)
		data := append(opt.Info, make([]string, 10)...)
		if val, err := strconv.Atoi(data[0]); err == nil {
			if val <= 100000 {
				// check if contest
				opt.contClass = "contest"
			} else {
				// type is gym
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
	// set path to folder containing contest
	opt.dirPath = path
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

/*
Parsing structure of problems
-----------------------------
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
