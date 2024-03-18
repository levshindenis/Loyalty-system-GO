package database

import (
	"context"
	"database/sql"
	"time"
)

func (dbs *DBStorage) CheckUserBalance(db *sql.DB, userID string, orderID string, orderSum float64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var balance float64
	err := db.QueryRowContext(ctx,
		`SELECT balance FROM balances WHERE user_id = $1`, userID).Scan(&balance)
	if err != nil {
		return false, err
	}

	if balance < orderSum {
		return false, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (order_id, user_id, status, accrual, uploaded_at) values ($1, $2, $3, $4, $5)`,
		orderID, userID, "SOLD", orderSum, time.Now().Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		return false, err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE balances SET balance = balance - $1, withdrawn = withdrawn + $1 WHERE user_id = $2`,
		orderSum, userID)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	return true, tx.Commit()
}
