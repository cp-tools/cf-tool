package cfg

import (
	pkg "cf/packages"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/infixint943/cookiejar"
)

// Session holds cookies and other request header data
// password is encrypted and stored securely
var Session struct {
	Handle  string         `json:"handle"`
	Passwd  string         `json:"password"`
	Cookies *cookiejar.Jar `json:"cookies"`
	Client  http.Client    `json:"-"`
}

var sessPath string

// InitSession reads data from sessions.json
func InitSession(path string) {
	// set sessions.json file path
	sessPath = path

	Session.Handle = ""
	Session.Cookies, _ = cookiejar.New(nil)
	proxyURL := http.ProxyFromEnvironment

	file, err := ioutil.ReadFile(sessPath)
	if err != nil {
		pkg.Log.Warning("File sessions.json doesn't exist")
		pkg.Log.Info("Creating sessions.json file...")
		SaveSession()
	}
	json.Unmarshal(file, &Session)
	// configure proxy if set
	if Settings.Proxy != "" {
		proxy, _ := url.Parse(Settings.Proxy)
		proxyURL = http.ProxyURL(proxy)
	}

	// instantiate client with proxy configurations
	Session.Client = http.Client{Jar: Session.Cookies,
		Transport: &http.Transport{Proxy: proxyURL}}
}

// SaveSession saves the data to sessions.json
func SaveSession() {
	// create sessions.json file and log err (if any)
	file, err := os.Create(sessPath)
	pkg.PrintError(err, "Failed to create sessions.json file")

	body, _ := json.MarshalIndent(Session, "", "\t")
	file.Write(body)
}
