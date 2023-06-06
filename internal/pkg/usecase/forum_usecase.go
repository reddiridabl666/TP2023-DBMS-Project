package usecase

import (
	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ForumUsecase struct {
	forums *repository.ForumRepository
	users  *repository.UserRepository
}

func NewForumUsecase(db *pgxpool.Pool) *ForumUsecase {
	return &ForumUsecase{
		forums: repository.NewForumRepository(db),
		users:  repository.NewUserRepository(db),
	}
}

func (u *ForumUsecase) Create(forum *domain.Forum) error {
	user, err := u.users.GetByNickname(forum.User)
	if err != nil {
		return err
	}

	forum.User = user.Nickname
	return u.forums.Create(user.Id, forum)
}

func (u *ForumUsecase) Get(slug string) (*domain.Forum, error) {
	return u.forums.Get(slug)
}
