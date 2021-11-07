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
	season   db.SeasonRepository
	config   config.TorrentServerConfig
	torrent  TorrentService
	telegram *TelegramService
}

func NewScraperService(repo db.SeasonRepository, conf config.TorrentServerConfig, torrent TorrentService, telegram *TelegramService) ScraperService {
	return ScraperService{
		season:   repo,
		config:   conf,
		torrent:  torrent,
		telegram: telegram,
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

			if ss.telegram != nil {
				telegramMessage := fmt.Sprintf("Added %s s%02de%02d",
					s.Name,
					s.Season,
					s.LastEpisode+1)

				if err := ss.telegram.SendMessage(telegramMessage); err != nil {
					log.Printf("Failed to send notification to telegram with: %v\n", err)
				}
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
		return getMagnetLinkFromAnchor(searchResult)
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
	if anchor != nil {
		magnetLink, ok = anchor.Attr("href")
	} else {
		return "", fmt.Errorf("did not find requested episode")
	}
	if ok {
		return magnetLink, nil
	} else {
		return "", fmt.Errorf("error extracting href tag")
	}
}
