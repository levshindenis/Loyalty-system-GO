package models

import "fmt"

type Task struct {
	OrderId string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual"`
}

func NewTask(orderId string, status string, accrual *float64) Task {
	return Task{
		OrderId: orderId,
		Status:  status,
		Accrual: accrual,
	}
}

func (t *Task) String() {
	fmt.Println(t.OrderId + "  " + t.Status)
}
