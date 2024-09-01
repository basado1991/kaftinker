package types

type User struct {
	Id int64 `json:"id"`

	Username string `json:"username"`
	Password []byte `json:"-"`
}
