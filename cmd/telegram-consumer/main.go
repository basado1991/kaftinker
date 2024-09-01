package main

import (
	"context"
	"log"
	"os"
	"strconv"

	telegramconsumer "example.org/nn/kaftinker/internal/telegram-consumer"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/segmentio/kafka-go"
)

const GROUP_ID = "telegram"

func getVar(name string) string {
	res, exists := os.LookupEnv(name)
	if !exists {
		log.Println(name, "is not set")
	}
	return res
}

func main() {
	kafkaAddr := getVar("KT_KAFKA_ADDR")
	kafkaTopic := getVar("KT_KAFKA_TOPIC")

	tgToken := getVar("TGC_TOKEN")
	tgChatId, err := strconv.ParseInt(getVar("TGC_CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Fatalln(err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddr},
		GroupID: GROUP_ID,
		Topic:   kafkaTopic,
	})

	ctx := context.Background()

	consumer := telegramconsumer.TelegramConsumer{
		Ctx:         ctx,
		KafkaReader: reader,

		TelegramApi:     bot,
		TelegramGroupId: tgChatId,
	}

	log.Println("starting telegram consumer...")
	consumer.Main()
}
