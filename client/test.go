package cln

import (
	cfg "cf/config"
	pkg "cf/packages"

	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
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
		err := errors.New("Unequal number of input/output test files")
		return nil, nil, err
	} else if len(inp) == 0 {
		err := errors.New("No test files found")
		return nil, nil, err
	}
	return inp, out, nil
}

// FindSourceFiles finds all code files in current dir
// with file name matching pattern
func FindSourceFiles(pattern string) (string, error) {
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
	// validate and set source file
	if len(files) == 0 {
		err := errors.New("No source files found\n" +
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
	pkg.PrintError(err, "")
	return file, nil
}

// FindTmpltsConfig finds all templates matching extension
// of `file` and returns all suitable template alias
func FindTmpltsConfig(file string) (*cfg.Template, error) {
	// index of valid template configurations
	var id []int
	for i, t := range cfg.Templates {
		// file extensions match = valid config
		if t.Ext == filepath.Ext(file) {
			id = append(id, i)
		}
	}
	// validate and set template config
	if len(id) == 0 {
		err := errors.New("No template configuration found\n" +
			"Ensure a suitable configured template exists")
		return nil, err
	} else if len(id) == 1 {
		// set template configuration (only 1 present)
		return &cfg.Templates[id[0]], nil
	}

	// prompt user for template configuration to select
	var idx int
	err := survey.AskOne(&survey.Select{
		Message: "Template configuration:",
		Options: cfg.ListTmplts(id...),
	}, &idx)
	pkg.PrintError(err, "")
	return &cfg.Templates[idx], nil
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
		// omit exp difference <= 1e-<exp>
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

// PrintDiff is run if outputs don't match
// prints input data, and then the diff of out vs ans
// prints all data to stderr, since it's debugging info
func PrintDiff(inp, out, ans string) {
	headerfmt := color.New(color.FgBlue, color.Bold, color.Underline).SprintfFunc()
	// print input data
	fmt.Fprintln(os.Stderr, headerfmt("Input"))
	fmt.Fprintln(os.Stderr, inp)

	// break output into lines
	str1 := strings.Split(out, "\n")
	str2 := strings.Split(ans, "\n")
	// equalize string lengths
	if len(str1) < len(str2) {
		str1 = append(str1, make([]string, len(str2)-len(str1))...)
	} else {
		str2 = append(str2, make([]string, len(str1)-len(str2))...)
	}

	// print output diff data
	tbl := uitable.New()
	tbl.Separator = " | "

	tbl.AddRow(headerfmt("Actual Output"), headerfmt("Expected Output"))
	// iterate over every row of outputs
	for i := 0; i < len(str1); i++ {
		tbl.AddRow(str1[i], str2[i])
	}
	fmt.Fprintln(os.Stderr, tbl)
	fmt.Println()

	return
}
