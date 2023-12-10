package models

import (
	"container/list"
	"sync"
)

type Queue struct {
	queue *list.List
	mutex sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		queue: list.New(),
	}
}

func (cq *Queue) Enqueue(car *Car) {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	cq.queue.PushBack(car)
}

func (cq *Queue) Dequeue() *Car {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.queue.Len() == 0 {
		return nil
	}
	element := cq.queue.Front()
	cq.queue.Remove(element)
	return element.Value.(*Car)
}

func (cq *Queue) First() *Car {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.queue.Len() == 0 {
		return nil
	}
	element := cq.queue.Front()
	return element.Value.(*Car)
}

func (cq *Queue) Last() *Car {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	if cq.queue.Len() == 0 {
		return nil
	}
	element := cq.queue.Back()
	return element.Value.(*Car)
}

func (cq *Queue) Size() int {
	cq.mutex.Lock()
	defer cq.mutex.Unlock()
	return cq.queue.Len()
}
