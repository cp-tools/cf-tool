package cfg

import (
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

var (
	// Templates holds all configured templates of user
	Templates []Template
	tmpltPath string
)

// InitTemplates reads data from templates.json
func InitTemplates(path string) error {
	// set templates.json file path
	tmpltPath = path

	file, err := os.OpenFile(tmpltPath, os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(file)
	json.Unmarshal(body, &Templates)
	return nil

}

// SaveTemplates to settings.json file
func SaveTemplates() error {
	// create templates.json file
	file, err := os.OpenFile(tmpltPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	body, _ := json.MarshalIndent(Templates, "", "\t")
	file.Write(body)
	return nil
}

// ListTmplts returns an array of required template aliases
// basically, just extracts all template aliases of tmplt
func ListTmplts(tmplt ...Template) (opts []string) {
	for _, t := range tmplt {
		opts = append(opts, t.Alias)
	}
	return
}
