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
	TransmissionConfig    TransmissionConfig
	TelegramServiceConfig TelegramServiceConfig
	TmdbAPIConfig         TmdbAPIConfig
}

type DBConfig struct {
	Username string
	Password string
	Url      string
	Port     string
	URI      string
}

type TransmissionConfig struct {
	Url      string
	Port     uint16
	Username string
	Password string
}

type TelegramServiceConfig struct {
	ChatId   string
	BotToken string
	Enabled  bool
}

type ServerConfig struct {
	Port        string
	ApiBasePath string
}

type TmdbAPIConfig struct {
	ApiKey string
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
		TransmissionConfig: TransmissionConfig{
			Username: envConfig["TRANSMISSION_USERNAME"],
			Password: envConfig["TRANSMISSION_PASSWORD"],
			Url:      envConfig["TRANSMISSION_SERVER"],
			Port:     StringToUint16(envConfig["TRANSMISSION_PORT"]),
		},
		TelegramServiceConfig: TelegramServiceConfig{
			ChatId:   envConfig["TELEGRAM_CHAT_ID"],
			BotToken: envConfig["TELEGRAM_BOT_TOKEN"],
			Enabled:  stringToBoolWithFallback(envConfig["TELEGRAM_ENABLE_NOTIFICATIONS"], false),
		},
		TmdbAPIConfig: TmdbAPIConfig{
			ApiKey: envConfig["TMDB_API_KEY"],
		},
	}

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
	isProduction := getEnvBool("IS_PROD", false)
	if isProduction {
		configSource = "config/.env.prod"
	} else {
		configSource = "config/.env.local"
	}
	return configSource
}

func getEnvBool(value string, fallback bool) bool {
	envValue := os.Getenv(value)
	if len(envValue) == 0 {
		return fallback
	}
	return stringToBoolWithFallback(envValue, fallback)
}

func stringToBoolWithFallback(val string, fallback bool) bool {
	valueBool, err := strconv.ParseBool(val)
	if err != nil {
		log.Fatalf("failed to parse string(%s) to bool with: %+v", val, err)
		return fallback
	}
	return valueBool
}

func StringToUint16(value string) uint16 {
	uInt64, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		log.Fatalf("failed to convert string(%s) to uint16 with: %+v", value, err)
	}
	return uint16(uInt64)
}
