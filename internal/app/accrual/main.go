package accrual

import (
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
)

type CompareWorker struct {
	id      int
	queue1  chan models.Task
	queue2  chan models.Task
	address string
}

func NewCompareWorker(id int, queue1 chan models.Task, queue2 chan models.Task, address string) CompareWorker {
	return CompareWorker{
		id:      id,
		queue1:  queue1,
		queue2:  queue2,
		address: address,
	}
}

func (cw *CompareWorker) Loop(sugarLogger *zap.SugaredLogger) {
	for {
		var task models.Task

		client := resty.New()

		client.
			SetRetryCount(3).
			SetRetryWaitTime(60 * time.Second).
			SetRetryMaxWaitTime(60 * time.Second).
			SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
				return 0, errors.New("quota exceeded")
			})

		client.AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			},
		)

		t := <-cw.queue1

		resp, err := client.R().SetResult(&task).Get(cw.address + "/api/orders/" + t.OrderID)
		if err != nil {
			sugarLogger.Infoln(
				"time", time.Now().Format(time.RFC3339),
				"error", "Error with Get",
			)
			continue
		}

		if resp.StatusCode() == 204 {
			continue
		}
		if task.Status != t.Status {
			nTask := task
			if task.Status == "REGISTERED" {
				nTask = models.NewTask(task.OrderID, "PROCESSING", task.Accrual)
			}
			cw.queue2 <- nTask
		}
	}
}
