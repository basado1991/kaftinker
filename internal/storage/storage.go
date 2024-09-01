package storage

import (
	"context"

	"example.org/nn/kaftinker/internal/types"
	"example.org/nn/kaftinker/internal/types/dto"
)

type Storage interface {
	GetUserById(context.Context, int64) (*types.User, error)
	GetUserByUsername(context.Context, string) (*types.User, error)
	GetPostsWithAuthors(context.Context, int64, int64) ([]dto.PostWithUser, error)
	GetPostById(context.Context, int64) (*types.Post, error)

	CreateUser(context.Context, types.User) (int64, error)
	CreatePost(context.Context, types.Post) (int64, error)
}
