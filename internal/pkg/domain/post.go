package domain

import (
	"forum/internal/pkg/utils"
	"forum/internal/pkg/utils/nullable"
)

//easyjson:json
type Post struct {
	Id       int64
	Parent   nullable.Int64 `json:"parent,omitempty"`
	Author   string
	Message  string
	IsEdited bool `json:"isEdited,omitempty"`
	Forum    string
	Thread   int32
	Created  utils.Time
}

//easyjson:json
type PostBatch []*Post

const (
	SortFlat = iota
	SortTree
	SortParent
)

type PostListParams struct {
	ThreadId int
	Limit    int
	Since    int
	Sort     int
	Desc     bool
}
