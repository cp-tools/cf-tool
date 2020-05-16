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

// FindCountdown parses countdown (if exists) from countdown page
func FindCountdown(contest string, link url.URL) (int64, error) {
	// This implementation contains redirection prevention
	c := cfg.Session.Client
	c.CheckRedirect = pkg.RedirectCheck
	link.Path = path.Join(link.Path, "countdown")
	body, err := pkg.GetReqBody(&c, link.String())
	if err != nil {
		return 0, err
	} else if len(body) == 0 {
		// such page doesn't exist
		err = fmt.Errorf("Contest %v doesn't exist", contest)
		return 0, err
	}

	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := doc.Find("span.countdown").Text()

	var h, m, s int64
	fmt.Sscanf(val, "%d:%d:%d", &h, &m, &s)
	return h*3600 + m*60 + s, nil
}

// FetchProbs finds all problems present in the contest
func FetchProbs(contest string, link url.URL) ([]string, error) {
	// no need of modifying link as it already points to dashboard
	c := cfg.Session.Client
	body, err := pkg.GetReqBody(&c, link.String())
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

// FetchTests extracts test cases of the problem(s) in contest
// Returns 2d slice mapping to input and output
// If problem == "", fetch all problem test cases
// else, only fetch of given problem.
// fix for https://github.com/infixint943/cf/pull/2#issuecomment-626122011
func FetchTests(contest, problem string, link url.URL) ([][]string, [][]string, error) {

	c := cfg.Session.Client
	if problem == "" {
		// fetch from problems page
		link.Path = path.Join(link.Path, "problems")
	} else {
		// fetch from individual problem page
		link.Path = path.Join(link.Path, "problem", problem)
	}

	body, err := pkg.GetReqBody(&c, link.String())
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
