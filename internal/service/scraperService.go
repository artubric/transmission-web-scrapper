package service

import (
	"context"
	"fmt"
	"log"
	"time"
	"transmission-web-scrapper/internal/db"
	"transmission-web-scrapper/internal/service/scraper"
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
					telegramMessage = fmt.Sprintf("%s.%s",telegramMessage, "%0A✓ Season complete")
				}
				
				if err := ss.telegram.SendMessage(telegramMessage); err != nil {
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
	scraperFunc, ok := scraper.SourceMap[s.DataSource.SourceType]
	if !ok {
		return "", fmt.Errorf("unknown source type: %s", s.DataSource.SourceType)
	}

	return scraperFunc(s)
}

