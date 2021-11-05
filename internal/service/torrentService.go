package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"transmission-web-scrapper/config"
)

type addTorrentRequest struct {
	Filename    string `json:"Filename"`
	DownloadDir string `json:"DownloadDir"`
}

type addTorrentResponse struct {
	HashString string `json:"hashString"`
	Id         int64  `json:"id"`
	Name       string `json:"name"`
}

type TorrentService struct {
	conf config.TorrentServerConfig
}

func NewTorrentService(conf config.TorrentServerConfig) TorrentService {
	return TorrentService{
		conf: conf,
	}
}

func (ts TorrentService) Add(magnetLink string, downloadDir string) error {
	request := addTorrentRequest{
		Filename:    magnetLink,
		DownloadDir: downloadDir,
	}
	requestByte, err := json.Marshal(request)

	if err != nil {
		log.Printf("Failed to marshallData with: %+v", err)
		return err
	}

	resp, err := http.Post(ts.conf.AddTorrentURI, "application/json",
		bytes.NewBuffer(requestByte))

	if err != nil {
		log.Printf("Failed to POST with: %+v", err)
	}

	response := addTorrentResponse{}
	json.NewDecoder(resp.Body).Decode(&response)

	if response.HashString == "" {
		err := fmt.Errorf("failed to add request to torrent server")
		return err
	}

	log.Printf("added torrent: %+v\n", response.Name)

	return nil
}
