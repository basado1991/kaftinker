package storage

import (
	"context"
	"database/sql"
	"time"

	"example.org/nn/kaftinker/internal/types"
	"example.org/nn/kaftinker/internal/types/dto"
	_ "github.com/mattn/go-sqlite3"
)

const SQLITE_REQUEST_TIMEOUT = 2000 * time.Millisecond

type SqliteStorage struct {
	conn *sql.DB
}

func NewSqliteStorage(sqliteDatabasePath string) (*SqliteStorage, error) {
	db, err := sql.Open("sqlite3", sqliteDatabasePath)
	if err != nil {
		return nil, err
	}

	return &SqliteStorage{conn: db}, nil
}

func (s *SqliteStorage) GetUserById(ctx context.Context, userId int64) (*types.User, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLITE_REQUEST_TIMEOUT)
	defer cancel()

	row := s.conn.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE id = ?", userId)

	var u types.User
	if err := row.Scan(&u.Id, &u.Username, &u.Password); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *SqliteStorage) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLITE_REQUEST_TIMEOUT)
	defer cancel()

	row := s.conn.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username = ?", username)

	var u types.User
	if err := row.Scan(&u.Id, &u.Username, &u.Password); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *SqliteStorage) GetPostsWithAuthors(ctx context.Context, offset int64, count int64) ([]dto.PostWithUser, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLITE_REQUEST_TIMEOUT)
	defer cancel()

	rows, err := s.conn.QueryContext(ctx, "SELECT posts.id, posts.title, posts.body, users.id, users.username FROM posts, users WHERE posts.author_id = users.id LIMIT ? OFFSET ?", count, offset)
	if err != nil {
		return nil, err
	}

	var posts []dto.PostWithUser
	for rows.Next() {
		var v dto.PostWithUser

		if err := rows.Scan(&v.Post.Id, &v.Post.Title, &v.Post.Body, &v.User.Id, &v.User.Username); err != nil {
			return nil, err
		}
		v.Post.AuthorId = v.User.Id

		posts = append(posts, v)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *SqliteStorage) GetPostById(ctx context.Context, id int64) (*types.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLITE_REQUEST_TIMEOUT)
	defer cancel()

	row := s.conn.QueryRowContext(ctx, "SELECT id, author_id, title, body FROM posts WHERE id=?", id)

	var v types.Post
	if err := row.Scan(&v.Id, &v.AuthorId, &v.Title, &v.Body); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *SqliteStorage) CreatePost(ctx context.Context, post types.Post) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLITE_REQUEST_TIMEOUT)
	defer cancel()

	res, err := s.conn.ExecContext(ctx, "INSERT INTO posts (author_id, title, body) VALUES (?, ?, ?)",
		post.AuthorId, post.Title, post.Body)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

func (s *SqliteStorage) CreateUser(ctx context.Context, user types.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLITE_REQUEST_TIMEOUT)
	defer cancel()

	res, err := s.conn.ExecContext(ctx, "INSERT INTO users (username, password) VALUES (?, ?)",
		user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return id, nil
}
