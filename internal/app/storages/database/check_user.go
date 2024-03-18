package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/generators"
)

func (dbs *DBStorage) CheckUser(db *sql.DB, login string, password string, param string) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userID, pswrd string
	err := db.QueryRowContext(ctx,
		`SELECT user_id, password FROM users WHERE login = $1`, login).Scan(&userID, &pswrd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if param == "registration" {
				var counter int
				err = db.QueryRowContext(ctx, `SELECT count(*) FROM users`).Scan(&counter)
				if err != nil {
					return false, "", err
				}

				cookie, err := generators.GenerateCookie(counter + 1)
				if err != nil {
					return false, "", err
				}

				tx, err := db.Begin()
				if err != nil {
					return false, "", err
				}

				_, err = tx.ExecContext(ctx,
					`INSERT INTO users (user_id, login, password) values ($1, $2, $3)`,
					cookie, login, password)
				if err != nil {
					tx.Rollback()
					return false, "", err
				}

				_, err = tx.ExecContext(ctx,
					`INSERT INTO balances (user_id, balance, withdrawn) values ($1, $2, $3)`,
					cookie, 0, 0)
				if err != nil {
					tx.Rollback()
					return false, "", err
				}
				return false, cookie, tx.Commit()
			}
			return false, "", nil
		}
		return false, "", err
	}

	if param == "login" && pswrd != password {
		return false, "", nil
	}

	return true, userID, nil
}
