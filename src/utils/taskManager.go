package utils

import (
	"sync"
	"time"
)

type TaskManager[T comparable] struct {
	tasks sync.Map
}

func NewTaskManager[T comparable]() *TaskManager[T] {
	return &TaskManager[T]{}
}

func (tm *TaskManager[T]) NewTaskWithFail(delay time.Duration, arg T, task func(T), fail func(T)) bool {
	if fail == nil {
		return tm.NewTask(delay, arg, task)
	}

	_, loaded := tm.tasks.LoadOrStore(arg, struct{}{})
	if loaded {
		return false
	}


	time.AfterFunc(delay, func() {
		defer func() {
			defer tm.tasks.Delete(arg)
			r := recover()
			if r != nil {
				fail(arg)
			}
			
		}()
		task(arg)
	})
	return true
}

func (tm *TaskManager[T]) NewTask(delay time.Duration, arg T, task func(T)) bool {
	_, loaded := tm.tasks.LoadOrStore(arg, struct{}{})
	if loaded {
		return false
	}
	tm.tasks.Store(arg, struct{}{})
	time.AfterFunc(delay, func() {
		defer func() {
			tm.tasks.Delete(arg)
		}()
		task(arg)
	})
	return true
}

func (tm *TaskManager[T]) ExistTask(key T) bool {
	_, loaded := tm.tasks.Load(key)
	return loaded
}
