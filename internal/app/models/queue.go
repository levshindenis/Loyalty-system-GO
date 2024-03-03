package models

type Queue struct {
	ch chan *Task
}

func NewQueue() *Queue {
	return &Queue{
		ch: make(chan *Task, 1024),
	}
}

func (q *Queue) Push(t *Task) {
	q.ch <- t
}

func (q *Queue) PopWait() *Task {
	return <-q.ch
}
