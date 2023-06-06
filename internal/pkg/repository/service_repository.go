package repository

import (
	"database/sql"

	"forum/internal/pkg/domain"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (repo *ServiceRepository) Clear() error {
	_, err := repo.db.Exec("TRUNCATE users CASCADE")
	return err
}

func (repo *ServiceRepository) Status() (*domain.ServiceInfo, error) {
	res := &domain.ServiceInfo{}

	err := repo.db.QueryRow("SELECT count(*) FROM users").Scan(&res.User)
	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow(
		`SELECT count(*), COALESCE(sum(thread_num), 0), COALESCE(sum(post_num), 0) FROM forum`).
		Scan(&res.Forum, &res.Thread, &res.Post)
	if err != nil {
		return nil, err
	}

	return res, nil
}
