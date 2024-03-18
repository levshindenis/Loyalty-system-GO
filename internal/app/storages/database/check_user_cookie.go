package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func (dbs *DBStorage) CheckUserCookie(cookie string) (bool, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return false, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userID string
	err = db.QueryRowContext(ctx,
		`SELECT user_id FROM users WHERE user_id = $1`, cookie).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}