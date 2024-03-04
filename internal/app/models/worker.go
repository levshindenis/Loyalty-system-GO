package models

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type CompareWorker struct {
	id      int
	queue1  Queue
	queue2  Queue
	address string
	cond    *sync.Cond
}

func NewCompareWorker(id int, queue1 Queue, queue2 Queue, address string, cond *sync.Cond) CompareWorker {
	return CompareWorker{
		id:      id,
		queue1:  queue1,
		queue2:  queue2,
		address: address,
		cond:    cond,
	}
}

func (cw *CompareWorker) Loop() {
	for {
		var buf bytes.Buffer
		var task Task

		t := cw.queue1.Pop()

		resp, err := http.Get(cw.address + "/api/orders/" + t.OrderId)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode == 200 {
			if _, err = buf.ReadFrom(resp.Body); err != nil {
				panic(err)
			}

			if err = resp.Body.Close(); err != nil {
				panic(err)
			}

			if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
				panic(err)
			}

			if task.Status != t.Status {
				if task.Status == "REGISTERED" {
					cw.queue2.Push(NewTask(task.OrderId, "PROCESSING", task.Accrual))
				} else {
					cw.queue2.Push(task)
				}
			}
		} else if resp.StatusCode == 204 {
			continue
		} else {
			cw.cond.L.Lock()
			cw.cond.Wait()
			timer := time.NewTimer(60 * time.Second)

			select {
			case <-timer.C:
				cw.cond.Signal()
			}
			cw.cond.L.Unlock()
		}

	}
}
