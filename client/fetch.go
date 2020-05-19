package cln

import (
	cfg "cf/config"

	"bytes"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
FindCountdown parses and returns number of seconds remaining
before contest begins. Returns 0 if countdown has already ended.
Virtual contests (of the current user session) are supported too.

If countdown page doesn't exsit, returns error ErrContestNotExists
*/
func FindCountdown(contest string, link url.URL) (int64, error) {
	// This implementation contains redirection prevention
	c := cfg.Session.Client
	c.CheckRedirect = redirectCheck
	link.Path = path.Join(link.Path, "countdown")
	body, err := getReqBody(&c, link.String())
	if err != nil {
		return 0, err
	} else if len(body) == 0 {
		// such page doesn't exist
		return 0, ErrContestNotExists
	}

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := doc.Find("span.countdown").Text()

	var h, m, s int64
	fmt.Sscanf(val, "%d:%d:%d", &h, &m, &s)
	return h*3600 + m*60 + s, nil
}

/*
FetchProbs parses and returns problem code's of all problems in contest.
Problem codes are returned in their lowercase versions. For example,
A => a, F1 => f1, C2 => c2 etc.

If contest dashboard doesn't exist, returns error ErrContestNotExists
*/
func FetchProbs(contest string, link url.URL) ([]string, error) {
	// no need of modifying link as it already points to dashboard
	c := cfg.Session.Client
	body, err := getReqBody(&c, link.String())
	if err != nil {
		return nil, err
	}

	var probs []string
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find(".problems .id a").Each(func(_ int, s *goquery.Selection) {
		prob := strings.TrimSpace(s.Text())
		probs = append(probs, strings.ToLower(prob))
	})
	return probs, nil
}

/*
FetchTests parses test cases of problems in the contest and returns
the sample inputs/outputs as a 2d slice of strings.

If problem parameter is empty, returns test cases of ALL problems in contest.
Samples are fetched from the contest's 'complete problemset' page.

Otherwise, sample tests of only the specified problem is returned
Here, samples are fetched from the (individual) problem's page

If problems page doesn't exist, returns error ErrContestNotExists
*/
func FetchTests(contest, problem string, link url.URL) ([][]string, [][]string, error) {

	c := cfg.Session.Client
	if problem == "" {
		// fetch from problems page
		link.Path = path.Join(link.Path, "problems")
	} else {
		// fetch from individual problem page
		link.Path = path.Join(link.Path, "problem", problem)
	}

	body, err := getReqBody(&c, link.String())
	if err != nil {
		return nil, nil, err
	}

	// splInp will hold input of each problem
	// splOut maps to splInp with the output data
	var splInp, splOut [][]string
	// Iterate over every problem
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find(".sample-test").Each(func(_ int, prob *goquery.Selection) {

		// func to clean sample input/output text
		f := func(_ int, text *goquery.Selection) string {
			str, _ := text.Html()
			str = strings.ReplaceAll(str, "<br/>", "\n")
			return strings.TrimSpace(str) + "\n"
		}
		// iterate over all input fields
		splInp = append(splInp, prob.Find(".input pre").Map(f))
		// iterate over all output fields
		splOut = append(splOut, prob.Find(".output pre").Map(f))
	})
	return splInp, splOut, nil
}
