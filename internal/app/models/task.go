package models

type Task struct {
	orderId string
	status  string
}

func NewTask(orderId string, status string) *Task {
	return &Task{
		orderId: orderId,
		status:  status,
	}
}
