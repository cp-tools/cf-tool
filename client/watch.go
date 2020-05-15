package cln

import (
	cfg "cf/config"
	pkg "cf/packages"

	"bytes"
	"fmt"
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

// WatchSubmissions finds all submissions in contID that matches query string
// query = problem to fetch all submissions in a particular problem (should be uppercase)
// query = submitID to fetch submission of given submission id
func WatchSubmissions(group, contest, query string, link url.URL) ([]Submission, error) {
	// This implementation contains redirection prevention
	c := cfg.Session.Client
	c.CheckRedirect = pkg.RedirectCheck
	// fetch all submissions in contest
	link.Path = path.Join(link.Path, "my")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return nil, err
	} else if len(body) == 0 {
		// such page doesn't exist
		err = fmt.Errorf("Contest %v doesn't exist", contest)
		return nil, err
	}
	// to hold all submissions
	var data []Submission

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	sel := doc.Find("tr[data-submission-id]").Has(`a[href*="/` + strings.ToUpper(query) + `"]`)
	sel.Each(func(_ int, row *goquery.Selection) {

		// compress verdict and return color coded string
		clean := func(verdict string) string {
			verdict = strings.ReplaceAll(verdict, "Wrong answer", "WA")
			verdict = strings.ReplaceAll(verdict, "Time limit exceeded", "TLE")
			verdict = strings.ReplaceAll(verdict, "Memory limit exceeded", "TLE")

			switch {
			case strings.HasPrefix(verdict, "TLE"):
				return pkg.Yellow.Sprint(verdict)
			case strings.HasPrefix(verdict, "MLE"):
				return pkg.Red.Sprint(verdict)
			case strings.HasPrefix(verdict, "WA"):
				return pkg.Red.Sprint(verdict)
			case strings.HasPrefix(verdict, "Pretests passed"):
				return pkg.Green.Sprint(verdict)
			case strings.HasPrefix(verdict, "Accepted"):
				return pkg.Green.Sprint(verdict)
			default:
				return verdict
			}
		}

		data = append(data, Submission{
			ID:      pkg.GetText(row, ".id-cell"),
			When:    pkg.GetText(row.Find("td").First().Next(), "*"),
			Name:    pkg.GetText(row, "td[data-problemId]"),
			Lang:    pkg.GetText(row, "td:not([class])"),
			Waiting: pkg.GetAttr(row, ".status-cell", "waiting"),
			Verdict: clean(pkg.GetText(row, ".status-verdict-cell")),
			Time:    pkg.GetText(row, ".time-consumed-cell"),
			Memory:  pkg.GetText(row, ".memory-consumed-cell"),
		})
	})

	return data, nil
}

// WatchContest parses contest solved count status
func WatchContest(group, contest string, link url.URL) ([]Problem, error) {
	// This implementation contains redirection prevention
	c := cfg.Session.Client
	c.CheckRedirect = pkg.RedirectCheck
	// fetch contest dashboard page
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return nil, err
	} else if len(body) == 0 {
		// such page doesn't exist
		err = fmt.Errorf("Contest %v doesn't exist", contest)
		return nil, err
	}
	// to hold all problems in contest
	var data []Problem

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find(".problems tr").Has("td").Each(func(_ int, row *goquery.Selection) {

		data = append(data, Problem{
			ID:     pkg.GetText(row, ".id"),
			Name:   pkg.GetText(row, "td > div > div > a"),
			Status: row.AttrOr("class", ""),
			Count:  pkg.GetText(row, "td > a"),
		})
	})
	return data, nil
}
