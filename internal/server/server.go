package server

import (
	"transmission-web-scrapper/config"
	"transmission-web-scrapper/internal/db"
)

type Server struct {
	config  config.ServerConfig
	dbRepos *db.DBRepositories
}

func New(conf config.ServerConfig, dbRepos *db.DBRepositories) Server {
	return Server{
		config:  conf,
		dbRepos: dbRepos,
	}
}

func (s Server) Run() {
	// init routes
}
