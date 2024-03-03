package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type CompareWorker struct {
	id      int
	queue1  *Queue
	queue2  *Queue
	address string
}

func NewCompareWorker(id int, queue1 *Queue, queue2 *Queue, address string) *CompareWorker {
	return &CompareWorker{
		id:      id,
		queue1:  queue1,
		queue2:  queue2,
		address: address,
	}
}

func (cw *CompareWorker) Loop() error {
	for {
		var buf bytes.Buffer
		var accOrder AccrualOrder

		t := cw.queue1.PopWait()

		resp, err := http.Get(cw.address + t.orderId)
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			if _, err = buf.ReadFrom(resp.Body); err != nil {
				return err
			}

			if err = resp.Body.Close(); err != nil {
				return err
			}

			if err = json.Unmarshal(buf.Bytes(), &accOrder); err != nil {
				return err
			}

			if accOrder.Status != t.status {
				cw.queue2.Push(NewTask(accOrder.OrderID, accOrder.Status))
			}
		} else {
			return errors.New("error from Worker Loop")
		}

	}
}
