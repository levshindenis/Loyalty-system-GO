package server

import (
	"context"
	"time"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

func (serv *Storage) FromChannelToDB(q chan models.Task, ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)

	var values []models.Task

	for {
		select {
		case <-ctx.Done():
			q <- models.Task{}
			return
		case value := <-q:
			values = append(values, value)
		case <-ticker.C:
			if len(values) == 0 {
				continue
			}
			err := serv.dbs.GetData().UpdateOrders(values)
			if err != nil {
				serv.sl.Infoln(
					"time", time.Now(),
					"error", "Error with FromChannelToDB",
				)
				continue
			}
			values = nil
		}
	}
}
