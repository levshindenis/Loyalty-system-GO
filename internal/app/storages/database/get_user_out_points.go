package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (dbs *DBStorage) GetUserOutPoints(db *sql.DB, userID string) (bool, []models.OutPoints, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx,
		`SELECT order_id, accrual, uploaded_at FROM orders 
            WHERE user_id = $1 and status = 'SOLD' order by uploaded_at desc`,
		userID)
	if err != nil {
		return false, nil, err
	}

	var items []models.OutPoints
	for rows.Next() {
		var item models.OutPoints
		if err = rows.Scan(&item.OrderID, &item.Summ, &item.ProcessedAt); err != nil {
			return false, nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return false, nil, err
	}

	if len(items) > 0 {
		return true, items, nil
	}

	return false, nil, nil
}
