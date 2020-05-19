package cln

import (
	cfg "cf/config"

	"encoding/hex"
	"net/url"
	"path"

	"github.com/infixint943/cookiejar"
	"github.com/oleiade/serrure/aes"
)

/*
Login attempts logging in to configured host domain
with user credentials passed in the parameters.

Returns true if login was successful (saves session to sessPath)
and false if login failed due to wrong credentials.

If login failed for any other reason (other than wrong creds)
the respective hhtp error message is returned.
*/
func Login(usr, passwd string) (bool, error) {
	// instantiate http client, but remove
	// past user sessions to prevent redirection
	jar, _ := cookiejar.New(nil)
	c := cfg.Session.Client
	c.Jar = jar

	link, _ := url.Parse(cfg.Settings.Host)
	link.Path = path.Join(link.Path, "enter")
	body, err := getReqBody(&c, link.String())
	if err != nil {
		return false, err
	}

	// Hidden form data
	csrf := findCsrf(body)
	ftaa := "yzo0kk4bhlbaw83g2q"
	bfaa := "883b704dbe5c70e1e61de4d8aff2da32"

	// Post form (aka login using creds)
	body, err = postReqBody(&c, link.String(), url.Values{
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

	usr = findHandle(body)
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

/*
LoggedInUsr returns handle of currently logged in user
Session uses Session.Client data to pull homepage
and extract the handle of the logged in user.
Returns an empty string if no logged in user is found

If http request failed, corresponding error is returned
*/
func LoggedInUsr() (string, error) {
	// fetch home page and check if logged in
	c := cfg.Session.Client
	link, _ := url.Parse(cfg.Settings.Host)
	body, err := getReqBody(&c, link.String())
	if err != nil {
		return "", err
	}

	return findHandle(body), nil
}

/*
Relogin extracts user handle / passwd from the Session struct
and passes the credentials to function Login() to relogin again.
Returns same return values of function Login()

If password couldn't be decrypted, returns error ErrDecodePasswdFailed
*/
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
