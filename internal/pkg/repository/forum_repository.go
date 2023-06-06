package repository

import (
	"context"
	"errors"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ForumRepository struct {
	db *pgxpool.Pool
}

func NewForumRepository(db *pgxpool.Pool) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (repo *ForumRepository) Create(userId int, forum *domain.Forum) error {
	_, err := repo.db.Exec(context.Background(),
		`INSERT INTO forum(title, slug, author_id)
		 	VALUES($1, $2, $3)`,
		forum.Title, forum.Slug, userId,
	)

	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != pgerrcode.UniqueViolation {
		return err
	}

	err = repo.db.QueryRow(context.Background(),
		`SELECT f.id, u.nickname, f.title,
				   f.slug, f.post_num, f.thread_num
			FROM forum f JOIN users u ON f.author_id = u.id
			WHERE lower(f.slug) = lower($1)`, forum.Slug).
		Scan(
			&forum.Id,
			&forum.User,
			&forum.Title,
			&forum.Slug,
			&forum.Posts,
			&forum.Threads,
		)
	if err != nil {
		return err
	}
	return domain.ErrAlreadyExists
}

func (repo *ForumRepository) Get(slug string) (*domain.Forum, error) {
	forum := &domain.Forum{}

	err := repo.db.QueryRow(context.Background(),
		`SELECT f.id, u.nickname, f.title,
			f.slug, f.post_num, f.thread_num
		 FROM forum f JOIN users u ON f.author_id = u.id
		 WHERE lower(f.slug) = lower($1)`, slug).
		Scan(
			&forum.Id,
			&forum.User,
			&forum.Title,
			&forum.Slug,
			&forum.Posts,
			&forum.Threads,
		)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return forum, nil
}
