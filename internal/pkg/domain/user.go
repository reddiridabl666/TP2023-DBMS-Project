package domain

import (
	"forum/internal/pkg/utils/nullable"
)

//easyjson:json
type User struct {
	Id       int `json:"-"`
	Nickname string
	Fullname string
	About    nullable.String
	Email    string
}

//easyjson:json
type UserBatch []*User

type UserListParams struct {
	ForumId int
	Limit   int
	Since   string
	Desc    bool
}
