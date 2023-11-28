package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"transmission-web-scrapper/internal/db"

	"github.com/PuerkitoBio/goquery"
)

var SourceMap = map[string]func(s db.Season) (string, error) {
	//
	"eztv" : func(s db.Season) (string, error){
		url := fmt.Sprintf("%s%s-%s-s%02de%02d",
			s.DataSource.Link,
			s.Name,
			s.Quality,
			s.Season,
			s.LastEpisode+1)
		log.Printf("Fetching html via URL: %s\n", url)
		document, err := getHTMLpage(url)
		if err != nil {
			return "", err
		}
		var searchResult *goquery.Selection
		numberOfElements := document.Find("tr.forum_header_border").Length()
		// eztv feature: if search found no torrents, it returns latest 50 torrents on first page
		if numberOfElements == 50 {
			searchResult = nil
		} else {
			searchResult = document.Find("a.magnet").First()
		}

		return getMagnetLinkFromAnchor(searchResult)
	},
	//
	"nyaa.si": func(s db.Season) (string, error){
		url := fmt.Sprintf("%s%s+%s+s%02de%02d%s",
			s.DataSource.Link,
			s.Name,
			s.Quality,
			s.Season,
			s.LastEpisode+1,
			s.DataSource.Parameters)
		log.Printf("Fetching html via URL: %s\n", url)

		document, err := getHTMLpage(url)
		if err != nil {
			return "", err
		}
		var foundAnchor *goquery.Selection
		var foundRow *goquery.Selection
		searchRows := document.Find("tr.default")
		searchRows.EachWithBreak(func(i int, searchRow *goquery.Selection) bool {
			searchRow.Find("td").EachWithBreak(func(i int, searchColumn *goquery.Selection) bool {
				colspan, found := searchColumn.Attr("colspan")
				if found && colspan == "2" {
					if !hasBadKodec(searchColumn.Text()) {
						foundRow = searchRow
						return false
					}
				}
				return true
			})
			if foundRow != nil {
				foundRow.Find("td.text-center").First().Find("a").EachWithBreak(func(i int, anchor *goquery.Selection) bool {
					href, found := anchor.Attr("href")
					if found && strings.Contains(href, "magnet:") {
						foundAnchor = anchor
						return false
					}
					return true
				})
			}
			if foundAnchor != nil {
				return false
			}
			return true
		})

		return getMagnetLinkFromAnchor(foundAnchor)
	},
	//
	"torrentz2.nz": func(s db.Season) (string, error){
		url := fmt.Sprintf("%s%s-%s-s%02de%02d",
			s.DataSource.Link,
			s.Name,
			s.Quality,
			s.Season,
			s.LastEpisode+1)

		log.Printf("Fetching html via URL: %s\n", url)
		
		document, err := getHTMLpage(url)
		if err != nil {
			return "", err
		}
		
		// to avoid fake torrents; expected real one to be aploaded to multiple sources
		resultsFound := document.Find("div.results dl").Length()
		if resultsFound < 5 {
			return "", fmt.Errorf("found less than 5 result, cutting off")
		}

		anchor := document.Find("div.results dl dd span a").First()

		href, found := anchor.Attr("href")
		if !found || !strings.Contains(href, "magnet:") {
			return "", fmt.Errorf("did not find requested episode")
		}
		
		return getMagnetLinkFromAnchor(anchor)
	},
}

func getHTMLpage(url string) (*goquery.Document, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func getMagnetLinkFromAnchor(anchor *goquery.Selection) (string, error) {
	var magnetLink string
	var ok bool

	if anchor == nil {
		return "", fmt.Errorf("did not find requested episode")
	}

	magnetLink, ok = anchor.Attr("href"); if !ok {
		return "", fmt.Errorf("error extracting href tag")
	}
	
	return magnetLink, nil
}

func hasBadKodec(torrentTitle string) bool {
	kodecsToIgnore := []string{"AV1"}
	for _, badKodec := range kodecsToIgnore {
		if strings.Contains(torrentTitle, badKodec) {
			return true
		}
	}

	return false
}