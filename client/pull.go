package cln

import (
	cfg "cf/config"
	pkg "cf/packages"

	"bytes"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

type (
	// Sub struct holds submission data
	// parsed from cf API
	Sub struct {
		Contest, Problem, Lang,
		Sid string
	}
)

// FetchSubs pulls submissions matching criteria
func FetchSubs(contest, problem, handle string) ([]Sub, error) {

	c := cfg.Session.Client
	link, _ := url.Parse(cfg.Settings.Host)
	link.Path = path.Join(link.Path, "api", `user.status`)
	// add parameter to link
	q := link.Query()
	q.Set("handle", handle)
	link.RawQuery = q.Encode()

	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return nil, err
	}

	// check API response status
	status := gjson.GetBytes(body, "status").String()
	if status != "OK" {
		comm := gjson.GetBytes(body, "comment").String()
		return nil, fmt.Errorf(comm)
	}
	// is another submission to same problem considered
	isParsed := make(map[string]bool)
	var Subs []Sub

	// parsing of json text takes place here.
	// thanks to module tidwall/gjson for the awesome package
	result := gjson.GetBytes(body, "result")
	result.ForEach(func(key, value gjson.Result) bool {
		// check if result matches search criteria
		// extract submission data
		contID := value.Get("problem.contestId").String()
		probID := value.Get("problem.index").String()
		probID = strings.ToLower(probID)

		verdict := value.Get("verdict").String()
		lang := value.Get("programmingLanguage").String()
		sid := value.Get("id").String()
		// ContestId+ProblemId => 1234c2
		query := contID + probID

		if (contID == contest || contest == "") && (probID == problem || problem == "") &&
			(verdict == "OK" && isParsed[query] == false) {
			// create sub and fetch source code
			s := Sub{
				Contest: contID,
				Problem: probID,
				Lang:    lang,
				Sid:     sid,
			}
			// push submission into struct
			Subs = append(Subs, s)
			// set to true, to prevent parsing other submissions of this problem
			isParsed[query] = true
		}
		// continue iteration
		return true
	})
	return Subs, nil
}

// FetchSubSource fetches submission code of Sub
func (sub *Sub) FetchSubSource() (string, error) {
	// determine contest type (contest/gym)
	contClass := "contest"
	if val, _ := strconv.Atoi(sub.Contest); val > 100000 {
		contClass = "gym"
	}

	c := cfg.Session.Client
	link, _ := url.Parse(cfg.Settings.Host)
	link.Path = path.Join(link.Path, contClass, sub.Contest, "submission", sub.Sid)
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return "", err
	}

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	source := doc.Find("pre#program-source-text").Text()
	return source, nil
}
