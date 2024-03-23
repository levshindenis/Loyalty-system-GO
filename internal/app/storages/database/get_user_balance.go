package database

import (
	"context"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (dbs *DBStorage) GetUserBalance(userID string) (models.Balance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var item models.Balance
	err := dbs.DB.QueryRowContext(ctx,
		`SELECT balance, withdrawn FROM balances WHERE user_id = $1`, userID).Scan(&item.Current, &item.WithDrawn)
	if err != nil {
		return models.Balance{}, err
	}
	return item, nil
}
