package service

import (
	"context"
	"log"
	"transmission-web-scrapper/config"

	trpc "github.com/hekmon/transmissionrpc/v2"
)

type TransmissionService struct {
	client *trpc.Client
}

func NewTransmissionService(conf config.TransmissionConfig) *TransmissionService {
	advancedConfig := &trpc.AdvancedConfig{
		Port: conf.Port,
	}
	log.Println("Connecting to transmission remote client... ")
	client, _ := trpc.New(
		conf.Url,
		conf.Username,
		conf.Password,
		advancedConfig,
	)

	// use as sanity, since err on client creation is empty
	ok, _, _, err := client.RPCVersion(context.Background())
	if !ok {
		log.Fatalf("failed to connect to transmission RPC:\n%+v", err)
	}

	ts := &TransmissionService{
		client: client,
	}

	return ts
}

func (ts TransmissionService) AddTorrent(ctx context.Context, downloadDir string, filename string) (trpc.Torrent, error) {
	torrent, err := ts.client.TorrentAdd(ctx,
		trpc.TorrentAddPayload{
			DownloadDir: &downloadDir,
			Filename:    &filename,
		},
	)

	return torrent, err
}
