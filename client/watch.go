package cln

import (
	cfg "cf/config"

	"bytes"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type (
	// Submission holds various data of
	// a particular submission
	Submission struct {
		ID, When, Name, Lang, Waiting,
		Verdict, Time, Memory string
	}
	// Problem holds problem solved status
	// based on current user session
	Problem struct {
		ID, Name, Status,
		Count string
	}
)

// WatchSubmissions finds all submissions in contest that matches query string
// query = problem to fetch all submissions in a particular problem (should be uppercase)
// query = submitID to fetch submission of given submission id
func WatchSubmissions(contest, query string, link url.URL) ([]Submission, error) {
	// This implementation contains redirection prevention
	c := cfg.Session.Client
	c.CheckRedirect = redirectCheck
	// fetch all submissions in contest
	link.Path = path.Join(link.Path, "my")
	body, err := getReqBody(&c, link.String())
	if err != nil {
		return nil, err
	} else if len(body) == 0 {
		// such page doesn't exist
		return nil, ErrContestNotExists
	}
	// to hold all submissions
	var data []Submission

	query = strings.ToUpper(query)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	sel := doc.Find("tr[data-submission-id]").Has("a[href*=\"/" + query + "\"]")
	sel.Each(func(_ int, row *goquery.Selection) {
		// select cell ...type(x) from row
		data = append(data, Submission{
			ID:      getText(row, "td:nth-of-type(1)"),
			When:    getText(row, "td:nth-of-type(2)"),
			Name:    getText(row, "td:nth-of-type(4)"),
			Lang:    getText(row, "td:nth-of-type(5)"),
			Waiting: getAttr(row, "td:nth-of-type(6)", "waiting"),
			Verdict: getText(row, "td:nth-of-type(6)"),
			Time:    getText(row, "td:nth-of-type(7)"),
			Memory:  getText(row, "td:nth-of-type(8)"),
		})
	})

	return data, nil
}

// WatchContest parses contest solved count status
func WatchContest(contest string, link url.URL) ([]Problem, error) {
	// This implementation contains redirection prevention
	c := cfg.Session.Client
	c.CheckRedirect = redirectCheck
	// fetch contest dashboard page
	body, err := getReqBody(&c, link.String())
	if err != nil {
		return nil, err
	} else if len(body) == 0 {
		// such page doesn't exist
		return nil, ErrContestNotExists
	}
	// to hold all problems in contest
	var data []Problem

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find(".problems tr").Has("td").Each(func(_ int, row *goquery.Selection) {

		data = append(data, Problem{
			ID:     getText(row, "td:nth-of-type(1)"),
			Name:   getText(row, "td:nth-of-type(2) a"),
			Count:  getText(row, "td:nth-of-type(4)"),
			Status: row.AttrOr("class", ""),
		})
	})
	return data, nil
}
