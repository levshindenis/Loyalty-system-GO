package database

import (
	"context"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (dbs *DBStorage) GetUserOrders(userID string) (bool, []models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := dbs.DB.QueryContext(ctx,
		`SELECT order_id, status, accrual, uploaded_at FROM orders 
            WHERE user_id = $1 and status <> 'SOLD' order by uploaded_at desc`,
		userID)
	if err != nil {
		return false, nil, err
	}

	var items []models.Order
	for rows.Next() {
		var item models.Order
		if err = rows.Scan(&item.OrderID, &item.Status, &item.Accrual, &item.UploadedAt); err != nil {
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
