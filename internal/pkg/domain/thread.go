package domain

import (
	"time"

	"forum/internal/pkg/utils"
	"forum/internal/pkg/utils/nullable"
)

//easyjson:json
type Thread struct {
	Id      int32
	Title   string
	Author  string
	Forum   string
	Message string
	Votes   int32           `json:"votes,omitempty"`
	Slug    nullable.String `json:"slug,omitempty"`
	Created utils.Time
}

//easyjson:json
type ThreadBatch []*Thread

type ThreadListParams struct {
	Forum   string
	ForumId int
	Limit   int
	Since   time.Time
	Desc    bool
}
