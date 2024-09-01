package types

type Post struct {
	Id int64 `json:"id"`

	AuthorId int64 `json:"author_id"`

	Title string `json:"title"`
	Body  string `json:"body"`
}
