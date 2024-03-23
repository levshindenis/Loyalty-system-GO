package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func (dbs *DBStorage) CheckUserOrder(orderID string, userID string) (bool, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var uid string
	err := dbs.DB.QueryRowContext(ctx,
		`SELECT user_id FROM orders WHERE order_id = $1`, orderID).Scan(&uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = dbs.DB.ExecContext(ctx,
				`INSERT INTO orders (order_id, user_id, status, accrual, uploaded_at) values ($1, $2, $3, $4, $5)`,
				orderID, userID, "NEW", nil, time.Now().Format(time.RFC3339))
			if err != nil {
				return false, false, err
			}
			uid = userID
			return false, false, nil
		}
		return false, false, err
	}

	if uid == userID {
		return true, true, nil
	}
	return true, false, nil
}
