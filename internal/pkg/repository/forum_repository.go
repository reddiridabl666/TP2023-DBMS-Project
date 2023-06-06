package repository

import (
	"database/sql"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) *ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (repo *ForumRepository) Create(forum *domain.Forum) error {
	_, err := repo.db.Exec(
		`INSERT INTO forum(title, slug, author_id)
		 VALUES($1, $2, (SELECT id FROM users WHERE lower(nickname) = lower($3)))`,
		forum.Title, forum.Slug, forum.User,
	)

	if err == nil {
		return nil
	}

	pgError, ok := err.(pgx.PgError)
	if !ok {
		return err
	}

	switch pgError.Code {
	case pgerrcode.UniqueViolation:
		err := repo.db.QueryRow(`
			SELECT f.id, u.nickname, f.title,
				   f.slug, f.post_num, f.thread_num
			FROM forum f JOIN users u ON f.author_id = u.id
			WHERE f.slug = $1`, forum.Slug).
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
		return domain.ErrUniqueViolation
	case pgerrcode.ForeignKeyViolation:
		return domain.ErrNotFound
	default:
		return err
	}
}

func (repo *ForumRepository) Get(slug string) (*domain.Forum, error) {
	forum := &domain.Forum{}

	err := repo.db.QueryRow(
		`SELECT f.id, u.nickname, f.title,
			f.slug, f.post_num, f.thread_num
		 FROM forum f JOIN users u ON f.author_id = u.id
		 WHERE f.slug = $1`, slug).
		Scan(
			&forum.Id,
			&forum.User,
			&forum.Title,
			&forum.Slug,
			&forum.Posts,
			&forum.Threads,
		)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return forum, nil
}