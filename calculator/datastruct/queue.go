package datastruct

import (
	"container/list"
)

type Queue[T any] struct {
	l *list.List
}

type QueueEle[T any] struct {
	v T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{l: list.New()}
}

func (q *Queue[T]) Enqueue(e T) {
	q.l.PushBack(&QueueEle[T]{e})
}

func (q *Queue[T]) Dequeue() T {
	var ret T
	e := q.dequeueEle()
	if e != nil {
		ret = e.v
	}

	return ret
}

func (q *Queue[T]) dequeueEle() *QueueEle[T] {
	if q.l.Len() == 0 {
		return nil
	}

	elem, ok := q.l.Remove(q.l.Front()).(*QueueEle[T])
	if !ok {
		return nil
	}

	return elem
}
