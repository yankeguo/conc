package conc

import (
	"context"
	"time"
)

// Periodic returns a Task do input Task periodically
func Periodic(task Task, dur time.Duration, init bool) Task {
	return TaskFunc(func(ctx context.Context) (err error) {
		if init {
			if err = task.Do(ctx); err != nil {
				return
			}
		}
		for {
			select {
			case <-time.After(dur):
			case <-ctx.Done():
				return
			}
			if err = task.Do(ctx); err != nil {
				return
			}
		}
	})
}
