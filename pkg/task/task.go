package task

import (
	"context"
	"fmt"
)

type TaskStatus int

const (
	New TaskStatus = iota
	Running
	Completed
)

var taskStatus = [3]string{"New", "Running", "Completed"}

type taskFunc[T interface{}] func(ctx context.Context) (T, error)

type task[T interface{}] struct {
	result T
	err    error
	c      chan struct{}
	status TaskStatus
	fn     taskFunc[T]
}

type TaskPanicError struct {
	reason interface{}
}

func (e TaskPanicError) Error() string {
	return fmt.Sprintf("Task panicked: %v", e.reason)
}

func NewTask[T interface{}](fn taskFunc[T]) *task[T] {
	t := new(task[T])
	t.c = make(chan struct{})
	t.fn = fn

	return t
}

func (t *task[T]) Status() TaskStatus {
	return t.status
}

func (t *task[T]) Result() (T, error) {
	switch t.status {
	case Completed:
		return t.result, t.err
	default:
		return t.result, fmt.Errorf("Task status is %s", taskStatus[t.status])
	}
}

func (t *task[T]) Await(ctx ...context.Context) (T, error) {
	switch t.status {
	case Completed:
		return t.result, t.err
	case Running:
		return t.result, fmt.Errorf("Task already running, can only be started once")
	}

	defer func() {
		t.status = Completed
	}()

	var _ctx context.Context
	if len(ctx) == 0 {
		_ctx = context.Background()
	} else {
		_ctx = ctx[0]
	}

	go t.start(_ctx)

	select {
	case <-_ctx.Done():
		return t.result, _ctx.Err()
	case <-t.c:
		return t.result, t.err
	}
}

func (t *task[T]) start(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			t.err = TaskPanicError{reason: r}
		}

		t.status = Completed
		t.c <- struct{}{}
		close(t.c)
	}()

	if ctx.Err() != nil {
		t.err = ctx.Err()
		t.c <- struct{}{}
		return
	}

	t.status = Running
	t.result, t.err = t.fn(ctx)
}
