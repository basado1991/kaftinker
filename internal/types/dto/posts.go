package dto

import "example.org/nn/kaftinker/internal/types"

type PostWithUser struct {
	types.Post
	types.User
}

type PostCreatedMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
