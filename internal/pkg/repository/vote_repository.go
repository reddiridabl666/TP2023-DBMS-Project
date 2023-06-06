package repository

import "github.com/jackc/pgx/v5/pgxpool"

type VoteRepository struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{
		db: db,
	}
}
