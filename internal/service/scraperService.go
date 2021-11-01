package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"transmission-web-scrapper/config"
	"transmission-web-scrapper/internal/db"

	"github.com/PuerkitoBio/goquery"
)

type ScraperService struct {
	season  db.SeasonRepository
	config  config.TorrentServerConfig
	torrent TorrentService
}

func NewScraperService(repo db.SeasonRepository, conf config.TorrentServerConfig, torrent TorrentService) ScraperService {
	return ScraperService{
		season:  repo,
		config:  conf,
		torrent: torrent,
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
			if err = ss.torrent.Add(magnetLink, s.DownloadDir); err != nil {
				log.Println(err)
				continue
			}

			if s.LastEpisode+1 >= s.TotalEpisodes {
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
		var searchResult *goquery.Selection
		anchorTags := document.Find("td.text-center").First().Find("a")
		anchorTags.EachWithBreak(func(i int, anchor *goquery.Selection) bool {
			href, found := anchor.Attr("href")
			if found && strings.Contains(href, "magnet:") {
				searchResult = anchor
				return false
			}
			return true
		})
		var magnetLink string
		var ok bool
		if searchResult != nil {
			magnetLink, ok = searchResult.Attr("href")
		} else {
			return "", fmt.Errorf("did not find requested episode")
		}
		if ok {
			return magnetLink, nil
		} else {
			return "", fmt.Errorf("error extracting href tag")
		}
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
