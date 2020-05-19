package cln

import (
	cfg "cf/config"

	"bytes"
	"context"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FindTests finds all returns all sample input/output
// files present in the current directory
func FindTests() ([]string, []string, error) {
	// function to read and return data of
	// all files matching pattern in dir
	readTests := func(pattern string) ([]string, error) {
		var data []string
		glob, _ := filepath.Glob(pattern)
		for _, file := range glob {
			str, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, err
			}
			data = append(data, string(str))
		}
		return data, nil
	}
	// fetch test cases
	inp, err := readTests("[0-9].in")
	if err != nil {
		return nil, nil, err
	}
	out, err := readTests("[0-9].out")
	if err != nil {
		return nil, nil, err
	}
	// validate input output data
	// check for i/o count equality
	// and existence of non-zero test files
	if len(inp) != len(out) {
		return nil, nil, ErrUnequalSampleTests
	} else if len(inp) == 0 {
		return nil, nil, ErrSampleTestsNotExists
	}
	return inp, out, nil
}

// FindSourceFiles finds all code files in current dir
// with file name matching pattern
func FindSourceFiles(pattern string) []string {
	// current pattern implementation follows *.*
	// and input/output files are excluded while checking
	// for existence of template config (L58-L66)
	glob, _ := filepath.Glob(pattern)
	var files []string
	// remove files not matching template extension
	for _, file := range glob {
		for _, t := range cfg.Templates {
			if t.Ext == filepath.Ext(file) {
				// insert as valid match found
				files = append(files, file)
				break
			}
		}
	}
	return files
}

// FindTmpltsConfig finds all templates matching extension
// of `file` and returns all suitable template alias
func FindTmpltsConfig(file string) []cfg.Template {
	// index of valid template configurations
	var templ []cfg.Template
	for _, t := range cfg.Templates {
		// file extensions match = valid config
		if t.Ext == filepath.Ext(file) {
			templ = append(templ, t)
		}
	}
	return templ
}

// ExecScript runs script with input and timeout and returns the
// time taken, stdout, stderr. Returns deadlineExceeded if timout occurs
func ExecScript(script, input string, dur int) (time.Duration, string, error) {
	cmds := strings.Split(script, " ")
	var stdout bytes.Buffer

	// set timer of `dur` seconds for execution of script
	secs := time.Duration(dur) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), secs)
	defer cancel()

	cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = io.Writer(&stdout)
	cmd.Stderr = os.Stderr

	// run script and measure time taken
	start := time.Now()
	err := cmd.Run()
	since := time.Since(start).Truncate(time.Millisecond)

	return since, stdout.String(), err
}

// Validator modifies and returns output / expected output
// based on flags passed (ignore-case / ignore-exp)
func Validator(out, ans string, igCase bool, exp int) (string, string) {
	// cleans the data based on validator flags
	f := func(data string) string {
		// remove trailing and leading spaces
		data = strings.TrimSpace(data)
		// convert to lower case if igCase set
		if igCase == true {
			data = strings.ToLower(data)
		}
		outData := ""
		// omit exp difference < 1e-<exp>
		for _, line := range strings.Split(data, "\n") {
			outLine := ""
			for _, wrd := range strings.Split(line, " ") {
				fVal, err := strconv.ParseFloat(wrd, 64)
				// valid float number, modify accordingly
				if err == nil {
					// truncate till exp places after point
					wrd = big.NewFloat(fVal).Text('f', exp)
					wrd = strings.TrimRight(strings.TrimRight(wrd, "0"), ".")
				}
				// join words with space in between
				outLine += wrd + " "
			}
			// remove trailing space
			outLine = strings.TrimSpace(outLine)
			// join lines with newline in between
			outData += outLine + "\n"
		}
		// remove trailing newline
		outData = strings.TrimSpace(outData)
		return outData
	}
	// return formatted strings
	return f(out), f(ans)
}
