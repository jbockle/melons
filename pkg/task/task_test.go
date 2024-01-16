package task

import (
	"context"
	"strings"
	"testing"
)

func TestNewTask(t *testing.T) {
	t.Run("NewTask", func(t *testing.T) {
		var taskFunc taskFunc[int] = func(ctx context.Context) (int, error) {
			return 1, nil
		}
		task := NewTask(taskFunc)

		if task == nil {
			t.Errorf("NewTask() = %v, want %v", task, "task")
		}
	})
}

func TestTaskAwait(t *testing.T) {
	t.Run("TaskAwait", func(t *testing.T) {
		task := NewTask(func(ctx context.Context) (int, error) {
			return 1, nil
		})

		result, err := task.Await()

		if err != nil {
			t.Errorf("TaskAwait() = %v, want %v", err, nil)
		}

		if result != 1 {
			t.Errorf("TaskAwait() = %v, want %v", result, 1)
		}
	})

	t.Run("TaskAwaitStatus", func(t *testing.T) {
		task := NewTask(func(ctx context.Context) (int, error) {
			return 1, nil
		})

		task.Await()

		if task.status != Completed {
			t.Errorf("TaskAwaitStatus() = %v, want %v", task.status, Completed)
		}
	})

	t.Run("TaskAlreadyStarted", func(t *testing.T) {
		task := NewTask(func(ctx context.Context) (int, error) {
			return 1, nil
		})
		task.status = Running

		_, err := task.Await()

		if err == nil {
			t.Error("Expected task already running error")
		}

		if !strings.Contains(err.Error(), "already running") {
			t.Errorf("TaskAlreadyStarted() = %v, want %v", err, "already running")
		}
	})
}
