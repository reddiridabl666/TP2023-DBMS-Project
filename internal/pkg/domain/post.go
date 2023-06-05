package domain

import "time"

//easyjson:json
type Post struct {
	Id       int64
	Parent   int64
	Author   string
	Message  string
	IsEdited bool
	Forum    string
	Thread   int32
	Created  time.Time
}
