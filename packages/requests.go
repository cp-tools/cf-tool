package pkg

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func parseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// GetReqBody executes a GET request to url and returns the request body
func GetReqBody(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return parseBody(resp)
}

// PostReqBody executes a POST request (with values: data) to url and returns the request body
func PostReqBody(client *http.Client, url string, data url.Values) ([]byte, error) {
	resp, err := client.PostForm(url, data)
	if err != nil {
		return nil, err
	}
	return parseBody(resp)
}

// FindHandle scrapes handle from REQUEST body
func FindHandle(body []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := doc.Find("#header").Find("a[href^=\"/profile/\"]").Text()
	return val
}

// FindCsrf extracts Csrf from REQUEST body
func FindCsrf(body []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val, _ := doc.Find(".csrf-token").Attr("data-csrf")
	return val
}
