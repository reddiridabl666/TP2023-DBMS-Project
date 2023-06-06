package delivery

import (
	"database/sql"

	"forum/internal/pkg/repository"
)

type ForumHandler struct {
	forums *repository.ForumRepository
}

func NewForumHandler(db *sql.DB) *ForumHandler {
	return &ForumHandler{
		forums: repository.NewForumRepository(db),
	}
}
