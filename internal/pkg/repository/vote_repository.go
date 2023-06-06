package repository

import (
	"context"
	"errors"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{
		db: db,
	}
}

func (repo *VoteRepository) Vote(vote *domain.Vote) error {
	_, err := repo.db.Exec(context.Background(),
		`INSERT INTO vote(author_id, thread_id, value)
			VALUES($1, $2, $3) ON CONFLICT (author_id, thread_id) 
			DO UPDATE SET value = EXCLUDED.value`,
		vote.UserId, vote.ThreadId, vote.Value,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
