package models

import "context"

type Queue struct {
	ch     chan Task
	ctx    context.Context
	cancel context.CancelFunc
}

func NewQueue() Queue {
	ch := make(chan Task, 1024)
	ctx, cancel := context.WithCancel(context.Background())
	return Queue{
		ch:     ch,
		ctx:    ctx,
		cancel: cancel}
}

func (q *Queue) GetCtx() context.Context {
	return q.ctx
}

func (q *Queue) Push(task Task) {
	q.ch <- task
}

func (q *Queue) GetChannel() chan Task {
	return q.ch
}

func (q *Queue) Pop() Task {
	return <-q.ch
}
