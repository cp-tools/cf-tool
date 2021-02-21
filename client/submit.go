package cln

import (
	cfg "cf/config"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
)

/*
Submit reads contents of file and submits it to
specified problem in contest. Returns nil is submission was successful.

If submission fails (includes failure due to submission of same code)
returns error message of cause of failed submission.
*/
func Submit(contest, problem, langID, file string, link url.URL) error {
	// form redirection prevention is removed while submitting
	c := cfg.Session.Client
	c.CheckRedirect = redirectCheck
	link.Path = path.Join(link.Path, "submit")
	body, err := getReqBody(&c, link.String())
	if err != nil {
		return err
	} else if len(body) == 0 {
		// such page doesn't exist
		return ErrContestNotExists
	}

	// read source file
	data, _ := ioutil.ReadFile(file)
	// hidden form data
	csrf := findCsrf(body)
	ftaa := "yzo0kk4bhlbaw83g2q"
	bfaa := "883b704dbe5c70e1e61de4d8aff2da32"
	// post form data (remove redirection prevention)
	c.CheckRedirect = nil
	body, err = postReqBody(&c, link.String(), url.Values{
		"csrf_token":            {csrf},
		"ftaa":                  {ftaa},
		"bfaa":                  {bfaa},
		"action":                {"submitSolutionFormSubmitted"},
		"submittedProblemIndex": {problem},
		"programTypeId":         {langID},
		"contestId":             {contest},
		"source":                {string(data)},
		"tabSize":               {"4"},
		"_tta":                  {"176"},
		"sourceCodeConfirmed":   {"true"},
	})
	if err != nil {
		return err
	}
	// find error message (if present)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	msg := doc.Find(".error").Text()
	if msg != "" {
		return fmt.Errorf(msg[2:])
	}

	return nil
}
