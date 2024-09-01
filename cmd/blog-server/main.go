package main

import (
	"context"
	"html/template"
	"log"
	"os"
	"strconv"

	blogserver "example.org/nn/kaftinker/internal/blog-server"
	"example.org/nn/kaftinker/internal/blog-server/handler"
	"example.org/nn/kaftinker/internal/blog-server/utils/cookie"
	"example.org/nn/kaftinker/internal/storage"
	"github.com/segmentio/kafka-go"
)

func getVar(name string) string {
	val, exists := os.LookupEnv(name)
	if !exists {
		log.Fatalln(name, "is not set")
	}

	return val
}

func main() {
	addr := getVar("BS_ADDR")
	assetsPath := getVar("BS_ASSETS_PATH")
	templatesPath := getVar("BS_TEMPLATES_PATH")
	publicKeyPath := getVar("BS_PUBLIC_KEY")
	privateKeyPath := getVar("BS_PRIVATE_KEY")
	sqliteDatabasePath := getVar("BS_SQLITE_DATABASE")
	passwordSalt := getVar("BS_PASSWORD_SALT")

	kafkaAddr := getVar("KT_KAFKA_ADDR")
	kafkaTopic := getVar("KT_KAFKA_TOPIC")
	kafkaPartition, err := strconv.ParseInt(getVar("KT_KAFKA_PARTITION"), 10, 0)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	kafkaConn, err := kafka.DialLeader(ctx, "tcp", kafkaAddr, kafkaTopic, int(kafkaPartition))
	if err != nil {
		log.Fatalln(err)
	}

	template, err := template.New("main").ParseGlob(templatesPath + "/*/*.html")
	if err != nil {
		log.Fatalln(err)
	}
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalln(err)
	}

	packer, err := cookie.NewCookiePacker(privateKey)
	if err != nil {
		log.Fatalln(err)
	}
	unpacker, err := cookie.NewCookieUnpacker(publicKey)
	if err != nil {
		log.Fatalln(err)
	}
	storage, err := storage.NewSqliteStorage(sqliteDatabasePath)
	if err != nil {
		log.Fatalln(err)
	}

	h := handler.Handler{
		Template:       template,
		CookiePacker:   packer,
		CookieUnpacker: unpacker,
		Ctx:            ctx,
		Storage:        storage,
		PasswordSalt:   passwordSalt,
		KafkaConn:      kafkaConn,
	}
	h.SetupRoutes(assetsPath)

	log.Println("starting server...")
	if err := blogserver.Serve(addr, h); err != nil {
		log.Fatalln("serve failed:", err)
	}
}
