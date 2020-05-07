package pkg

import (
	"os"

	"github.com/fatih/color"
)

// Log is struct holding functions to print colored to stderr
// (lightweight replacement for Logger)
var (
	writer = os.Stderr
	green  = color.New(color.FgGreen)
	blue   = color.New(color.FgBlue)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)

	Log struct {
		Success, Notice, Info, Error,
		Warning func(text ...interface{})
	}
)

func init() {
	Log.Success = func(text ...interface{}) { green.Fprintln(writer, text...) }
	Log.Notice = func(text ...interface{}) { color.New().Fprintln(writer, text...) }
	Log.Info = func(text ...interface{}) { blue.Fprintln(writer, text...) }
	Log.Error = func(text ...interface{}) { red.Fprintln(writer, text...) }
	Log.Warning = func(text ...interface{}) { yellow.Fprintln(writer, text...) }
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
