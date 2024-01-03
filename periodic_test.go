package conc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPeriodic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var v int
	task := TaskFunc(func(ctx context.Context) error {
		v++
		return nil
	})
	pt := Periodic(task, time.Millisecond*200, true)
	go pt.Do(ctx)

	time.Sleep(time.Millisecond * 500)
	cancel()

	assert.Equal(t, 3, v)
}

func TestPeriodicWithError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var v int
	task := TaskFunc(func(ctx context.Context) error {
		v++
		return errors.New("failed")
	})
	var err error
	pt := Periodic(task, time.Millisecond*200, true)
	go func() { err = pt.Do(ctx) }()

	time.Sleep(time.Millisecond * 500)
	cancel()

	assert.Equal(t, 1, v)
	assert.Error(t, err)
}

func TestPeriodicWithError2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var v int
	task := TaskFunc(func(ctx context.Context) error {
		v++
		if v == 2 {
			return errors.New("failed")
		} else {
			return nil
		}
	})
	var err error
	pt := Periodic(task, time.Millisecond*200, true)
	go func() { err = pt.Do(ctx) }()

	time.Sleep(time.Millisecond * 500)
	cancel()

	assert.Equal(t, 2, v)
	assert.Error(t, err)
}
