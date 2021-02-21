package cfg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/infixint943/cookiejar"
)

var (
	// Session holds cookies and other request header data
	// password is AES encrypted and stored securely
	Session struct {
		Handle  string         `json:"handle"`
		Passwd  string         `json:"password"`
		Cookies *cookiejar.Jar `json:"cookies"`
		Client  http.Client    `json:"-"`
	}
	// path to sessions.json file
	sessPath string
)

// set default values of Session struct
func init() {
	Session.Handle = ""
	Session.Passwd = ""
	Session.Cookies, _ = cookiejar.New(nil)
	Session.Client = *http.DefaultClient
}

// InitSession reads data from sessions.json
func InitSession(path string) error {
	// set sessions.json file path
	sessPath = path

	file, err := os.OpenFile(sessPath, os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(file)
	json.Unmarshal(body, &Session)
	// proxy configuration
	proxyURL := http.ProxyFromEnvironment
	if Settings.Proxy != "" {
		proxy, _ := url.Parse(Settings.Proxy)
		proxyURL = http.ProxyURL(proxy)
	}

	// instantiate client with proxy configurations
	Session.Client = http.Client{Jar: Session.Cookies,
		Transport: &http.Transport{Proxy: proxyURL}}

	return nil
}

// SaveSession saves the data to sessions.json
func SaveSession() error {
	// create sessions.json file and log err (if any)
	file, err := os.OpenFile(sessPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	body, _ := json.MarshalIndent(Session, "", "\t")
	file.Write(body)
	return nil
}
