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
	allSeasons, err := ss.season.GetAllExpanded(ctx)
	if err != nil {
		return err
	}
	for _, s := range allSeasons {
		if !s.IsArchived &&
			now.After(s.StartDate.Time()) &&
			now.Before(s.EndDate.Time()) {
			magnetLink, err := doScraping(s)
			if err != nil {
				return err
			}
			if err = ss.torrent.Add(magnetLink, s.DownloadDir); err != nil {
				return err
			}

			if s.LastEpisode+1 >= s.TotalEpisodes {
				s.LastEpisode = s.TotalEpisodes
				s.IsArchived = true
			} else {
				s.LastEpisode++
			}

			updatedSeason := db.Season{
				ID:            s.ID,
				Name:          s.Name,
				Season:        s.Season,
				StartDate:     s.StartDate,
				EndDate:       s.EndDate,
				TotalEpisodes: s.TotalEpisodes,
				LastEpisode:   s.LastEpisode,
				Quality:       s.Quality,
				DataSource:    s.DataSource.ID,
				IsArchived:    s.IsArchived,
				DownloadDir:   s.DownloadDir,
			}
			_, err = ss.season.Update(ctx, updatedSeason)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func doScraping(s db.SeasonExpanded) (string, error) {
	switch s.DataSource.SourceType {
	case "nyaa.si":
		// TODO: generate string from config
		url := fmt.Sprintf("https://nyaa.si/?f=0&c=1_2&q=%s+%s+s%02de%02d&s=seeders&o=desc", s.Name, s.Quality, s.Season, s.LastEpisode+1)
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
