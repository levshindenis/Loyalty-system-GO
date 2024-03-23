package server

import (
	"context"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (serv *Storage) FromDBToChannel(q chan models.Task, ctx context.Context) {
	ticker := time.NewTicker(4 * time.Second)

	for {
		select {
		case <-ctx.Done():
			q <- models.Task{}
			return
		case <-ticker.C:
			items, err := serv.dbs.GetNewOrders()
			if err != nil {
				serv.sl.Infoln(
					"time", time.Now(),
					"error", "Error with FromDBChannel",
				)
				continue
			}
			for _, elem := range items {
				q <- elem
			}
		}
	}
}
