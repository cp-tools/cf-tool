package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	// Settings holds configured settings data of the tool
	Settings struct {
		DfltTmplt  int    `json:"default_template"`
		GenOnFetch bool   `json:"gen_on_fetch"`
		Host       string `json:"host"`
		Proxy      string `json:"proxy"`
		WSName     string `json:"workspace_name"`
	}

	settPath string
)

func init() {
	// initialise default values of Settings struct
	Settings.DfltTmplt = -1
	Settings.GenOnFetch = false
	Settings.Host = "https://codeforces.com"
	Settings.Proxy = ""
	Settings.WSName = "codeforces"
}

// InitSettings reads settings.json file
func InitSettings(path string) error {
	// set settings.json file path
	settPath = path

	file, err := os.OpenFile(settPath, os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(file)
	json.Unmarshal(body, &Settings)
	return nil
}

// SaveSettings to settings.json file
func SaveSettings() error {
	// create settings.json file
	file, err := os.OpenFile(settPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	body, _ := json.MarshalIndent(Settings, "", "\t")
	file.Write(body)
	return nil
}
