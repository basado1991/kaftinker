package telegramconsumer

import (
	"context"
	"encoding/json"
	"log"

	"example.org/nn/kaftinker/internal/types/dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/segmentio/kafka-go"
)

type TelegramConsumer struct {
	KafkaReader   *kafka.Reader
	TelegramToken string

	TelegramApi     *tgbotapi.BotAPI
	TelegramGroupId int64

	Ctx context.Context
}

func (c TelegramConsumer) SendMessage(post dto.PostCreatedMessage) error {
	msg := tgbotapi.NewMessage(c.TelegramGroupId, post.Title+"\n\n"+post.Body)
	_, err := c.TelegramApi.Send(msg)

	return err
}

func (c TelegramConsumer) Main() {
	for {
		m, err := c.KafkaReader.ReadMessage(c.Ctx)
		if err != nil {
			log.Println(err)
			continue
		}
		var post dto.PostCreatedMessage
		if err := json.Unmarshal(m.Value, &post); err != nil {
			log.Println(err)
			continue
		}

		if err := c.SendMessage(post); err != nil {
			log.Println(err)
		}
	}
}
