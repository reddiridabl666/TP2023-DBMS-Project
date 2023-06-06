package usecase

import (
	"database/sql"
	"strconv"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"
)

type ThreadUsecase struct {
	forums  *repository.ForumRepository
	threads *repository.ThreadRepository
	users   *repository.UserRepository
}

func NewThreadUsecase(db *sql.DB) *ThreadUsecase {
	return &ThreadUsecase{
		forums:  repository.NewForumRepository(db),
		threads: repository.NewThreadRepository(db),
		users:   repository.NewUserRepository(db),
	}
}

func (u *ThreadUsecase) Create(thread *domain.Thread) error {
	user, err := u.users.GetByNickname(thread.Author)
	if err != nil {
		return err
	}

	forum, err := u.forums.Get(thread.Forum)
	if err != nil {
		return err
	}

	thread.Author = user.Nickname
	thread.Forum = forum.Slug
	return u.threads.Create(forum.Id, user.Id, thread)
}

func (u *ThreadUsecase) Get(slugOrId string) (*domain.Thread, error) {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		return u.threads.GetById(id)
	}
	return u.threads.GetBySlug(slugOrId)
}
