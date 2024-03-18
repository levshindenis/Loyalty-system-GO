package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (dbs *DBStorage) UpdateOrders(db *sql.DB, tasks []models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		_, err = tx.ExecContext(ctx,
			`UPDATE orders SET status = $1, accrual = $2, uploaded_at = $3 WHERE order_id = $4`,
			task.Status, *task.Accrual, time.Now().Format(time.RFC3339), task.OrderID)
		if err != nil {
			tx.Rollback()
			return err
		}

		var userID string
		err = tx.QueryRowContext(ctx,
			`SELECT user_id FROM orders WHERE order_id = $1`, task.OrderID).Scan(&userID)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE balances SET balance = balance + $1 WHERE user_id = $2`, *task.Accrual, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
