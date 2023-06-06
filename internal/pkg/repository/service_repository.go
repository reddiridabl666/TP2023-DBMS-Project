package repository

import (
	"context"

	"forum/internal/pkg/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (repo *ServiceRepository) Clear() error {
	_, err := repo.db.Exec(context.Background(), "TRUNCATE users CASCADE")
	return err
}

func (repo *ServiceRepository) Status() (*domain.ServiceInfo, error) {
	res := &domain.ServiceInfo{}

	err := repo.db.QueryRow(context.Background(), "SELECT count(*) FROM users").Scan(&res.User)
	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow(context.Background(),
		`SELECT count(*), COALESCE(sum(thread_num), 0), COALESCE(sum(post_num), 0) FROM forum`).
		Scan(&res.Forum, &res.Thread, &res.Post)
	if err != nil {
		return nil, err
	}

	return res, nil
}
