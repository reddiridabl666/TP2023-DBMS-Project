package domain

import (
	"time"

	"forum/internal/pkg/utils/nullable"
)

//easyjson:json
type Thread struct {
	Id      int32
	Title   string
	Author  string
	Forum   string
	Message string
	Votes   int32
	Slug    nullable.String
	Created time.Time
}
