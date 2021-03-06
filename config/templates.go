package cfg

import (
	pkg "cf/packages"

	"encoding/json"
	"io/ioutil"
	"os"
)

// Template struct manages templates data
type Template struct {
	LangName   string `json:"lang_name"`
	LangID     string `json:"lang_id"`
	Path       string `json:"path"`
	Ext        string `json:"ext"`
	Alias      string `json:"alias"`
	PreScript  string `json:"pre_script"`
	Script     string `json:"script"`
	PostScript string `json:"post_script"`
}

// Templates holds all configured templates of user
var Templates []Template

var tmpltPath string

// InitTemplates reads data from templates.json
func InitTemplates(path string) {
	// set templates.json file path
	tmpltPath = path

	file, err := ioutil.ReadFile(tmpltPath)
	if err != nil {
		pkg.Log.Warning("File templates.json doesn't exist")
		pkg.Log.Info("Creating templates.json file...")
		SaveTemplates()
	}
	json.Unmarshal(file, &Templates)
}

// SaveTemplates to settings.json file
func SaveTemplates() {
	file, err := os.Create(tmpltPath)
	pkg.PrintError(err, "Failed to create templates.json file")

	body, _ := json.MarshalIndent(Templates, "", "\t")
	file.Write(body)
}

// ListTmplts returns an array of required template aliases
// basically, just extracts all template aliases of tmplt
func ListTmplts(tmplt ...Template) (opts []string) {
	for _, t := range tmplt {
		opts = append(opts, t.Alias)
	}
	return
}
