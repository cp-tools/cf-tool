package cln

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Some global variables
var (
	ErrContestNotExists     = fmt.Errorf("Contest doesn't exist")
	ErrDecodePasswdFailed   = fmt.Errorf("Failed to decode password")
	ErrUnequalSampleTests   = fmt.Errorf("Unequal number of input/output test files")
	ErrSampleTestsNotExists = fmt.Errorf("No test files found")

	// LangID represents all available languages with id's
	LangID = map[string]string{
		"GNU GCC C11 5.1.0":                "43",
		"Clang++17 Diagnostics":            "52",
		"GNU G++11 5.1.0":                  "42",
		"GNU G++14 6.4.0":                  "50",
		"GNU G++17 7.3.0":                  "54",
		"Microsoft Visual C++ 2010":        "2",
		"Microsoft Visual C++ 2017":        "59",
		"GNU G++17 9.2.0 (64 bit, msys 2)": "61",
		"C# Mono 5.18":                     "9",
		"D DMD32 v2.086.0":                 "28",
		"Go 1.12.6":                        "32",
		"Haskell GHC 8.6.3":                "12",
		"Java 11.0.5":                      "60",
		"Java 1.8.0_162":                   "36",
		"Kotlin 1.3.10":                    "48",
		"OCaml 4.02.1":                     "19",
		"Delphi 7":                         "3",
		"Free Pascal 3.0.2":                "4",
		"PascalABC.NET 3.4.2":              "51",
		"Perl 5.20.1":                      "13",
		"PHP 7.2.13":                       "6",
		"Python 2.7.15":                    "7",
		"Python 3.7.2":                     "31",
		"PyPy 2.7 (7.2.0)":                 "40",
		"PyPy 3.6 (7.2.0)":                 "41",
		"Ruby 2.0.0p645":                   "8",
		"Rust 1.35.0":                      "49",
		"Scala 2.12.8":                     "20",
		"JavaScript V8 4.8.0":              "34",
		"Node.js 9.4.0":                    "55",
		"ActiveTcl 8.5":                    "14",
		"Io-2008-01-07 (Win32)":            "15",
		"Pike 7.8":                         "17",
		"Befunge":                          "18",
		"OpenCobol 1.0":                    "22",
		"Factor":                           "25",
		"Secret_171":                       "26",
		"Roco":                             "27",
		"Ada GNAT 4":                       "33",
		"Mysterious Language":              "38",
		"FALSE":                            "39",
		"Picat 0.9":                        "44",
		"GNU C++11 5 ZIP":                  "45",
		"Java 8 ZIP":                       "46",
		"J":                                "47",
		"Microsoft Q#":                     "56",
		"Text":                             "57",
	}

	// LangExt corresponds to file extension of
	// given language source code
	LangExt = map[string]string{
		"GNU C11":               ".c",
		"Clang++17 Diagnostics": ".cpp",
		"GNU C++0x":             ".cpp",
		"GNU C++":               ".cpp",
		"GNU C++11":             ".cpp",
		"GNU C++14":             ".cpp",
		"GNU C++17":             ".cpp",
		"MS C++":                ".cpp",
		"MS C++ 2017":           ".cpp",
		"GNU C++17 (64)":        ".cpp",
		"Mono C#":               ".cs",
		"D":                     ".d",
		"Go":                    ".go",
		"Haskell":               ".hs",
		"Kotlin":                ".kt",
		"Ocaml":                 ".ml",
		"Delphi":                ".pas",
		"FPC":                   ".pas",
		"PascalABC.NET":         ".pas",
		"Perl":                  ".pl",
		"PHP":                   ".php",
		"Python 2":              ".py",
		"Python 3":              ".py",
		"PyPy 2":                ".py",
		"PyPy 3":                ".py",
		"Ruby":                  ".rb",
		"Rust":                  ".rs",
		"JavaScript":            ".js",
		"Node.js":               ".js",
		"Q#":                    ".qs",
		"Java":                  ".java",
		"Java 6":                ".java",
		"Java 7":                ".java",
		"Java 8":                ".java",
		"Java 9":                ".java",
		"Java 10":               ".java",
		"Java 11":               ".java",
		"Tcl":                   ".tcl",
		"F#":                    ".fs",
		"Befunge":               ".bf",
		"Pike":                  ".pike",
		"Io":                    ".io",
		"Factor":                ".factor",
		"Cobol":                 ".cbl",
		"Secret_171":            ".secret_171",
		"Ada":                   ".adb",
		"FALSE":                 ".f",
		"":                      ".txt",
	}
)

func parseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// getReqBody executes a GET request to url and returns the request body
func getReqBody(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return parseBody(resp)
}

// postReqBody executes a POST request (with values: data) to url and returns the request body
func postReqBody(client *http.Client, url string, data url.Values) ([]byte, error) {
	resp, err := client.PostForm(url, data)
	if err != nil {
		return nil, err
	}
	return parseBody(resp)
}

// findHandle scrapes handle from REQUEST body
func findHandle(body []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val := doc.Find("#header").Find("a[href^=\"/profile/\"]").Text()
	return val
}

// findCsrf extracts Csrf from REQUEST body
func findCsrf(body []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	val, _ := doc.Find(".csrf-token").Attr("data-csrf")
	return val
}

// redirectCheck prevents redirection and returns requested page info
func redirectCheck(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

// getText extracts text from particular html data
func getText(sel *goquery.Selection, query string) string {
	str := sel.Find(query).Text()
	return strings.TrimSpace(str)
}

// getAttr extracts attribute valur of particular html data
func getAttr(sel *goquery.Selection, query, attr string) string {
	str := sel.Find(query).AttrOr(attr, "")
	return strings.TrimSpace(str)
}
