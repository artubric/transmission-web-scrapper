package service

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"time"
	"transmission-web-scrapper/internal/db"

	"github.com/PuerkitoBio/goquery"
)

type ScraperService struct {
	season      db.SeasonRepository
	transmisson *TransmissionService
	telegram    *TelegramService
}

func NewScraperService(repo db.SeasonRepository, transmisson *TransmissionService, telegram *TelegramService) ScraperService {
	return ScraperService{
		season:      repo,
		transmisson: transmisson,
		telegram:    telegram,
	}
}

func (ss ScraperService) Start(ctx context.Context) error {
	now := time.Now()
	allSeasons, err := ss.season.GetAll(ctx, true)
	if err != nil {
		return err
	}
	for _, s := range allSeasons {
		if !s.IsArchived &&
			now.After(s.StartDate.Time()) &&
			now.Before(s.EndDate.Time()) {
			magnetLink, err := scrapForMagnetLink(s)
			if err != nil {
				log.Println(err)
				continue
			}
			_, err = ss.transmisson.AddTorrent(context.Background(), s.DownloadDir, magnetLink)
			if err != nil {
				log.Println(err)
				continue
			}

			isSeasonComplete := s.LastEpisode+1 >= s.TotalEpisodes

			if ss.telegram != nil {
				telegramMessage := fmt.Sprintf("≫ %s s%02de%02d",
					s.Name,
					s.Season,
					s.LastEpisode+1)

				if isSeasonComplete {
					telegramMessage = fmt.Sprintf("%s.%s",telegramMessage, "\n✓ Season complete")
				}
				
				if err := ss.telegram.SendMessage(html.EscapeString(telegramMessage)); err != nil {
					log.Printf("Failed to send notification to telegram with: %v\n", err)
				}
			}

			if isSeasonComplete {
				s.LastEpisode = s.TotalEpisodes
				s.IsArchived = true
			} else {
				s.LastEpisode++
			}
			_, err = ss.season.Update(ctx, s)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	return nil
}

func scrapForMagnetLink(s db.Season) (string, error) {
	switch s.DataSource.SourceType {
	case "nyaa.si":
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
	case "eztv":
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
	case "torrentz2.nz":
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

		anchor := document.Find("dl dd span a").First()

		href, found := anchor.Attr("href")
		if found && strings.Contains(href, "magnet:") {
			getMagnetLinkFromAnchor(anchor)
		}
		
		return getMagnetLinkFromAnchor(nil)
	default:
		return "", fmt.Errorf("unknown source type: %s", s.DataSource.SourceType)
	}

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
