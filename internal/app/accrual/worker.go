package accrual

import (
	"bytes"
	"encoding/json"
	"github.com/levshindenis/Loyalty-system-GO/internal/app/models"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type CompareWorker struct {
	id      int
	queue1  models.Queue
	queue2  models.Queue
	address string
	cond    *sync.Cond
}

func NewCompareWorker(id int, queue1 models.Queue, queue2 models.Queue, address string, cond *sync.Cond) CompareWorker {
	return CompareWorker{
		id:      id,
		queue1:  queue1,
		queue2:  queue2,
		address: address,
		cond:    cond,
	}
}

func (cw *CompareWorker) Loop(sugarLogger *zap.SugaredLogger) {
	for {
		var buf bytes.Buffer
		var task models.Task

		t := cw.queue1.Pop()

		resp, err := http.Get(cw.address + "/api/orders/" + t.OrderID)
		if err != nil {
			sugarLogger.Infoln(
				"time", time.Now().Format(time.RFC3339),
				"error", "Error with http.Get",
			)
			continue
		}

		if resp.StatusCode == 200 {
			if _, err = buf.ReadFrom(resp.Body); err != nil {
				sugarLogger.Infoln(
					"time", time.Now().Format(time.RFC3339),
					"error", "Error with read body",
				)
				continue
			}

			if err = resp.Body.Close(); err != nil {
				sugarLogger.Infoln(
					"time", time.Now().Format(time.RFC3339),
					"error", "Error with close body",
				)
				continue
			}

			if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
				sugarLogger.Infoln(
					"time", time.Now().Format(time.RFC3339),
					"error", "Error with Unmarshal",
				)
				continue
			}

			if task.Status != t.Status {
				if task.Status == "REGISTERED" {
					cw.queue2.Push(models.NewTask(task.OrderID, "PROCESSING", task.Accrual))
				} else {
					cw.queue2.Push(task)
				}
			}
		} else if resp.StatusCode == 204 {
			continue
		} else {
			// Пытаюсь заблокировать worker на минуту, чтобы подождать таймаут
			// (при большем, чем N количестве запросов, может вернуться, что превышен лимит. Пытаюсь переждать минуту)
			// Это скорее всего неправильно. Хотел бы получить уточнение по этому моменту.
			cw.cond.L.Lock()
			cw.cond.Wait()
			timer := time.NewTimer(60 * time.Second)
			<-timer.C
			cw.cond.Signal()
			cw.cond.L.Unlock()
		}

	}
}
