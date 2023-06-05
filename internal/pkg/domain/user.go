package domain

import (
	"forum/internal/pkg/utils/nullable"
)

//easyjson:json
type User struct {
	Id       int
	Nickname string
	Fullname string
	About    nullable.String
	Email    string
}
