package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"transmission-web-scrapper/config"
)

const (
	parseMode = "Markdown"
)

type telegramResponse struct {
	Ok               bool    `json:"ok"`
	ErrorDescription *string `json:"description"`
}

type TelegramService struct {
	conf config.TelegramServiceConfig
}

func NewTelegramService(conf config.TelegramServiceConfig) TelegramService {
	return TelegramService{
		conf: conf,
	}
}

func (ts *TelegramService) SendMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s=&parse_mode=%s&text=%s",
		ts.conf.BotToken,
		ts.conf.ChatId,
		parseMode,
		message,
	)
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	response := telegramResponse{}
	json.NewDecoder(resp.Body).Decode(&response)

	if !response.Ok {
		return fmt.Errorf(*response.ErrorDescription)
	}

	return nil

	// response examples:
	// {"ok":false,"error_code":400,"description":"Bad Request: message text is empty"}
	// {"ok":true,"result":{"message_id":4,"from":{"id":x,"is_bot":true,"first_name":"x","username":"x"},"chat":{"id":x,"title":"x","type":"group","all_members_are_administrators":true},"date":1636200104,"text":"Hello"}}%
}
