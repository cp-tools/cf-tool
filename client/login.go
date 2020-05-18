package cln

import (
	cfg "cf/config"
	pkg "cf/packages"

	"encoding/hex"
	"net/url"
	"path"

	"github.com/infixint943/cookiejar"
	"github.com/oleiade/serrure/aes"
)

// Login tries logginging in with user creds
func Login(usr, passwd string) (bool, error) {
	// instantiate http client, but remove
	// past user sessions to prevent redirection
	jar, _ := cookiejar.New(nil)
	c := cfg.Session.Client
	c.Jar = jar

	link, _ := url.Parse(cfg.Settings.Host)
	link.Path = path.Join(link.Path, "enter")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return false, err
	}

	// Hidden form data
	csrf := pkg.FindCsrf(body)
	ftaa := "yzo0kk4bhlbaw83g2q"
	bfaa := "883b704dbe5c70e1e61de4d8aff2da32"

	// Post form (aka login using creds)
	body, err = pkg.PostReqBody(&c, link.String(), url.Values{
		"csrf_token":    {csrf},
		"action":        {"enter"},
		"ftaa":          {ftaa},
		"bfaa":          {bfaa},
		"handleOrEmail": {usr},
		"password":      {passwd},
		"_tta":          {"176"},
		"remember":      {"on"},
	})
	if err != nil {
		return false, err
	}

	usr = pkg.FindHandle(body)
	if usr != "" {
		// create aes 256 encryption and encode as
		// hex string and save to sessions.json
		enc, _ := aes.NewAES256Encrypter(usr, nil)
		ed, _ := enc.Encrypt([]byte(passwd))
		ciphertext := hex.EncodeToString(ed)
		// update sessions data
		cfg.Session.Cookies = jar
		cfg.Session.Handle = usr
		cfg.Session.Passwd = ciphertext
		cfg.SaveSession()
	}
	return (usr != ""), nil
}

// LoggedInUsr checks and returns whether
// current session is logged in
func LoggedInUsr() (string, error) {
	// fetch home page and check if logged in
	c := cfg.Session.Client
	link, _ := url.Parse(cfg.Settings.Host)
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return "", err
	}

	return pkg.FindHandle(body), nil
}

// Relogin extracts handle/passwd from sessions.json
// and log's in with the credentials and returns status
func Relogin() (bool, error) {
	// decode hex data of encrypted password
	ciphertext, err := hex.DecodeString(cfg.Session.Passwd)
	if err != nil {
		return false, ErrDecodePasswdFailed
	}
	usr := cfg.Session.Handle
	dec := aes.NewAES256Decrypter(usr)
	passwd, err := dec.Decrypt(ciphertext)
	if err != nil {
		return false, err
	}
	return Login(usr, string(passwd))
}
