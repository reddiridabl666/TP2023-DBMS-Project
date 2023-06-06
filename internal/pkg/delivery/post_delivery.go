package delivery

import (
	"database/sql"

	"forum/internal/pkg/repository"
)

type PostHandler struct {
	posts *repository.PostRepository
}

func NewPostHandler(db *sql.DB) *PostHandler {
	return &PostHandler{
		posts: repository.NewPostRepository(db),
	}
}
