package common

import (
	"github.com/go-redis/redis"
)

// TaskScheduler schedules tasks to be worked on.
type TaskScheduler interface {
	ScheduleTask(Task) error
	GetTaskResults() ([]Task, error)
}

type taskScheduler struct {
	client *redis.Client
}

// NewTaskScheduler creates a new TaskScheduler.
func NewTaskScheduler(redisURL string) TaskScheduler {
	c := redis.NewClient(&redis.Options{Addr: redisURL})
	return &taskScheduler{client: c}
}

func (ts *taskScheduler) ScheduleTask(task Task) error {
	s, err := SerializeTask(task)
	if err != nil {
		return err
	}

	err = ts.client.LPush(PendingTaskQueue, s).Err()
	if err != nil {
		return err
	}

	// Prevent too many tasks from building up.
	// In a real system, the number of workers should probably be scaled up.
	// This is an O(1) operation (since the worst case is always removing 1).
	err = ts.client.LTrim(PendingTaskQueue, 0, EnvConfig.MaxTaskQueueSize-1).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ts *taskScheduler) GetTaskResults() ([]Task, error) {
	return nil, nil
}
