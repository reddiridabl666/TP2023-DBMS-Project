package usecase

import (
	"database/sql"

	"forum/internal/pkg/repository"
)

type PostUsecase struct {
	posts   *repository.PostRepository
	threads *repository.ThreadRepository
}

func NewPostUsecase(db *sql.DB) *PostUsecase {
	return &PostUsecase{
		posts:   repository.NewPostRepository(db),
		threads: repository.NewThreadRepository(db),
	}
}
