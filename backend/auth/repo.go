package auth

import (
	"database/sql"
	"fmt"
)

func wrapTransaction(db *sql.DB, transaction func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = transaction(tx)
	if err != nil {
		err = fmt.Errorf("sql transaction failed: %v", err)
		return err
	}

	return tx.Commit()
}

type AuthRepo struct {
	db *sql.DB
}

func newAuthRepo(db *sql.DB) *AuthRepo {
	store := AuthRepo{
		db: db,
	}
	return &store
}

func (repo *AuthRepo) createUser(req *registerRequest, userId, secret string) error {
	err := wrapTransaction(repo.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			insert into users (id, secret, email, first, last)
			values (?, ?, ?, ?, ?);`,
			userId, secret, req.Email, req.First, req.Last)
		return err
	})
	return err
}

func (repo *AuthRepo) getUser(req *loginRequest) (*user, error) {
	row := repo.db.QueryRow(`select id, secret from users where email = ?;`, req.Email)

	u := user{}
	err := row.Scan(&u.id, &u.secret)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
