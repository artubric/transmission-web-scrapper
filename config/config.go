package config

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
}

type DBConfig struct {
	Username string
	Password string
	Url      string
	Port     string
}

type ServerConfig struct {
	Port string
}

func GetConfig() *Config {

	config := Config{
		DBConfig: DBConfig{
			Username: "zzz",
			Password: "zzz",
			Url:      "zzz",
			Port:     "zzz",
		},
		ServerConfig: ServerConfig{
			Port: "zzz",
		},
	}

	return &config
}
