package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig              DBConfig
	ServerConfig          ServerConfig
	TorrentServerConfig   TorrentServerConfig
	TelegramServiceConfig TelegramServiceConfig
}

type DBConfig struct {
	Username string
	Password string
	Url      string
	Port     string
	URI      string
}

type TorrentServerConfig struct {
	Protocol      string
	Address       string
	Port          string
	Endpoint      string
	AddTorrentURI string
}

type TelegramServiceConfig struct {
	ChatId   string
	BotToken string
}

type ServerConfig struct {
	Port        string
	ApiBasePath string
}

func Load() *Config {

	envConfig, err := godotenv.Read(getConfigSource())
	if err != nil {
		log.Fatalf("Failed to parse config from %+v", getConfigSource())
	}

	config := Config{
		DBConfig: DBConfig{
			Username: envConfig["DB_USERNAME"],
			Password: envConfig["DB_PASSWORD"],
			Url:      envConfig["DB_HOST"],
			Port:     envConfig["DB_PORT"],
		},
		ServerConfig: ServerConfig{
			Port:        envConfig["SERVER_PORT"],
			ApiBasePath: "api",
		},
		TorrentServerConfig: TorrentServerConfig{
			Protocol: envConfig["TORRENT_REST_PROTOCOL"],
			Address:  envConfig["TORRENT_REST_SERVER"],
			Port:     envConfig["TORRENT_REST_PORT"],
			Endpoint: envConfig["TORRENT_REST_ENDPOINT"],
		},
		TelegramServiceConfig: TelegramServiceConfig{
			ChatId:   envConfig["TELEGRAM_CHAT_ID"],
			BotToken: envConfig["TELEGRAM_BOT_TOKEN"],
		},
	}
	config.TorrentServerConfig.AddTorrentURI = fmt.Sprintf("%s://%s:%s/%s",
		config.TorrentServerConfig.Protocol,
		config.TorrentServerConfig.Address,
		config.TorrentServerConfig.Port,
		config.TorrentServerConfig.Endpoint,
	)

	config.DBConfig.URI = fmt.Sprintf("mongodb://%s:%s@%s:%s",
		config.DBConfig.Username,
		config.DBConfig.Password,
		config.DBConfig.Url,
		config.DBConfig.Port,
	)

	return &config
}

func getConfigSource() string {
	var configSource string
	isProduction := GetEnvBool("IS_PROD", false)
	if isProduction {
		configSource = "config/.env.prod"
	} else {
		configSource = "config/.env.local"
	}
	return configSource
}

func GetEnvBool(value string, fallback bool) bool {
	envValue := os.Getenv(value)
	if len(envValue) == 0 {
		return fallback
	}
	valueBool, err := strconv.ParseBool(envValue)
	if err != nil {
		log.Printf("failed to parse string(%s) to bool with: %+v", value, err)
		return fallback
	}
	return valueBool
}
