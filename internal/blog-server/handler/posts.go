package handler

import (
	"encoding/json"
	"log"
	"math"
	"net/http"

	"example.org/nn/kaftinker/internal/blog-server/utils"
	"example.org/nn/kaftinker/internal/types"
	"example.org/nn/kaftinker/internal/types/dto"
	"github.com/segmentio/kafka-go"
)

type PostCreationPageData struct {
	TitleNotProvided bool
	BodyNotProvided  bool
}

type IndexPageData struct {
	Entries []dto.PostWithUser
}

func (h Handler) getIndex(w http.ResponseWriter, r *http.Request) {
	_, err := utils.GetAuthenticatedUser(h.CookieUnpacker, r)
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}

	data, err := h.Storage.GetPostsWithAuthors(h.Ctx, 0, math.MaxInt64)
	if err != nil {
		log.Println(err)

		utils.WriteInternalError(w)
		return
	}

	pageData := IndexPageData{Entries: data}

	if err := h.Template.ExecuteTemplate(w, "index-page", pageData); err != nil {
		log.Println(err)

		return
	}
}

func (h Handler) getCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := h.Template.ExecuteTemplate(w, "create-post-page", PostCreationPageData{}); err != nil {
		log.Println(err)
	}
}

func (h Handler) postCreatePost(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetAuthenticatedUser(h.CookieUnpacker, r)
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		if err := h.Template.ExecuteTemplate(w, "create-post-page", PostCreationPageData{TitleNotProvided: true}); err != nil {
			log.Println(err)
		}
		return
	}
	body := r.FormValue("body")
	if body == "" {
		if err := h.Template.ExecuteTemplate(w, "create-post-page", PostCreationPageData{BodyNotProvided: true}); err != nil {
			log.Println(err)
		}
		return
	}

	newPost := types.Post{
		Title: title,
		Body:  body,

		AuthorId: userId,
	}

	_, err = h.Storage.CreatePost(h.Ctx, newPost)
	if err != nil {
		log.Println(err)

		utils.WriteInternalError(w)
		return
	}

	defer http.Redirect(w, r, "/", http.StatusFound)

	buff, err := json.Marshal(dto.PostCreatedMessage{Title: title, Body: body})
	if err != nil {
		log.Println(err)
		return
	}

	_, err = h.KafkaConn.WriteMessages(kafka.Message{Value: buff})
	if err != nil {
		log.Println("kafka error:", err)
	}
}
