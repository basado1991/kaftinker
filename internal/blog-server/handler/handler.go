package handler

import (
	"context"
	"html/template"
	"net/http"

	"example.org/nn/kaftinker/internal/blog-server/utils/cookie"
	"example.org/nn/kaftinker/internal/storage"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	Template *template.Template

	CookiePacker   *cookie.CookiePacker
	CookieUnpacker *cookie.CookieUnpacker

	Ctx          context.Context
	Storage      storage.Storage
	PasswordSalt string

	KafkaConn *kafka.Conn
}

func (h Handler) SetupRoutes(assetsPath string) {
	http.HandleFunc("GET /", h.getIndex)

	http.HandleFunc("GET /auth", h.getAuth)
	http.HandleFunc("POST /auth", h.postAuth)
	http.HandleFunc("GET /register", h.getRegister)
	http.HandleFunc("POST /register", h.postRegister)

	http.HandleFunc("GET /createPost", h.getCreatePost)
	http.HandleFunc("POST /createPost", h.postCreatePost)

	http.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath))))
}
