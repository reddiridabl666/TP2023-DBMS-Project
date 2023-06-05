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

// var tables = []string{"vote", "post", "thread", "forum", "users"}

func (repo *ServiceRepository) Clear() error {
	// return utils.Tx(repo.db, func(tx *sql.Tx) error {
	// 	for _, table := range tables {
	// 		_, err := tx.Exec("TRUNCATE " + table)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// })
	_, err := repo.db.Exec("TRUNCATE users CASCADE")
	return err
}

func (repo *ServiceRepository) Status() (*domain.ServiceInfo, error) {
	res := &domain.ServiceInfo{}

	err := repo.db.QueryRow("SELECT count(*) FROM users").Scan(&res.User)
	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow("SELECT count(*), sum(thread_num), sum(post_num) FROM forum").
		Scan(&res.Forum, &res.Thread, &res.Post)
	if err != nil {
		return nil, err
	}

	return res, nil
}
