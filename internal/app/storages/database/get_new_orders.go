package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (dbs *DBStorage) GetNewOrders() ([]models.Task, error) {
	db, err := sql.Open("pgx", dbs.GetAddress())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx,
		`SELECT order_id, status, accrual FROM orders WHERE status in ('NEW', 'PROCESSING')`)
	if err != nil {
		return nil, err
	}

	var items []models.Task
	for rows.Next() {
		var item models.Task
		if err = rows.Scan(&item.OrderID, &item.Status, &item.Accrual); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
