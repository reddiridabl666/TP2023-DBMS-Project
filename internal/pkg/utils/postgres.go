package utils

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

func Tx(db *sql.DB, fb func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fb(tx)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return err
	}

	return tx.Commit()
}

const (
	dsn      = "user=postgres dbname=forum password=12345 host=localhost port=5432 sslmode=disable"
	maxConns = 20
)

func InitPostgres() (*sql.DB, error) {
	till := time.Now().Add(time.Second * 10)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	for time.Now().Before(till) {
		log.Println("Trying to open pg connection")

		err = db.Ping()
		if err == nil {
			log.Println("Ping sucessful")
			break
		}

		time.Sleep(time.Second)
	}

	if err != nil {
		db.SetMaxOpenConns(maxConns)
	}
	return db, nil
}
