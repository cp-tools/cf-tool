package cln

import (
	cfg "cf/config"
	pkg "cf/packages"

	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
)

// Submit uploads form data and submits user code
func Submit(group, contest, problem, langID, file string, link url.URL) error {
	// form redirection prevention is removed while submitting
	c := cfg.Session.Client
	c.CheckRedirect = pkg.RedirectCheck
	link.Path = path.Join(link.Path, "submit")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return err
	} else if len(body) == 0 {
		// such page doesn't exist
		err = fmt.Errorf("Contest %v doesn't exist", contest)
		return err
	}

	// read source file
	data, _ := ioutil.ReadFile(file)
	// hidden form data
	csrf := pkg.FindCsrf(body)
	ftaa := "yzo0kk4bhlbaw83g2q"
	bfaa := "883b704dbe5c70e1e61de4d8aff2da32"
	// post form data (remove redirection prevention)
	c.CheckRedirect = nil
	body, err = pkg.PostReqBody(&c, link.String(), url.Values{
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
		return errors.New(msg[2:])
	}

	return nil
}
