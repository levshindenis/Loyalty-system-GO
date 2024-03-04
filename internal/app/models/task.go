package models

import "fmt"

type Task struct {
	OrderID string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual"`
}

func NewTask(orderID string, status string, accrual *float64) Task {
	return Task{
		OrderID: orderID,
		Status:  status,
		Accrual: accrual,
	}
}

func (t *Task) String() {
	fmt.Println(t.OrderID + "  " + t.Status)
}
