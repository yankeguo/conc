package conc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	var err = errors.New("test")
	var t1r, t2r, t3r bool
	var t1 TaskFunc = func(ctx context.Context) error {
		ta := time.After(time.Millisecond * 500)
		select {
		case <-ta:
			t1r = true
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		t2r = true
		return err
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		ta := time.After(time.Millisecond * 500)
		select {
		case <-ta:
			t3r = true
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	nt := Parallel(t1, t2, t3)
	ne := nt.Do(context.Background())
	assert.Error(t, ne)
	assert.False(t, t1r)
	assert.True(t, t2r)
	assert.False(t, t3r)

	errs := ne.(Errs)
	assert.Equal(t, context.Canceled, errs[0])
	assert.Equal(t, err, errs[1])
	assert.Equal(t, context.Canceled, errs[2])
}

func TestParallelWithLimit(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var err = errors.New("test")
	n := 100
	l := 10
	fn := 3
	trs := make([]bool, n, n)
	ts := make([]Task, n, n)

	var cc int64
	check := func() {
		assert.True(t, cc <= int64(l))
	}

	for _i := 0; _i < n; _i++ {
		i := _i
		ts[i] = TaskFunc(func(ctx context.Context) error {
			check()
			defer check()
			ta := time.After(time.Millisecond * 500)
			check()
			select {
			case <-ta:
				check()
				trs[i] = true
				check()
				return nil
			case <-ctx.Done():
				check()
				return ctx.Err()
			}
		})
	}

	ts[fn] = TaskFunc(func(ctx context.Context) error {
		trs[fn] = true
		return err
	})

	nt := ParallelWithLimit(l, ts...)
	ne := nt.Do(context.Background())
	assert.Error(t, ne)

	for i := 0; i < n; i++ {
		if i == fn {
			assert.True(t, trs[i])
		} else {
			assert.False(t, trs[i])
		}
	}

	errs := ne.(Errs)
	for i := 0; i < n; i++ {
		if i == fn {
			assert.Equal(t, err, errs[i])
		} else {
			assert.Equal(t, context.Canceled, errs[i])
		}
	}
}

func TestParallelFailSafe(t *testing.T) {
	var err = errors.New("test")
	var t1r, t2r, t3r bool
	var t1 TaskFunc = func(ctx context.Context) error {
		ta := time.After(time.Millisecond * 500)
		select {
		case <-ta:
			t1r = true
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		t2r = true
		return err
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		ta := time.After(time.Millisecond * 500)
		select {
		case <-ta:
			t3r = true
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	nt := ParallelFailSafe(t1, t2, t3)
	ne := nt.Do(context.Background())
	assert.Error(t, ne)
	assert.True(t, t1r)
	assert.True(t, t2r)
	assert.True(t, t3r)

	errs := ne.(Errs)
	assert.Equal(t, nil, errs[0])
	assert.Equal(t, err, errs[1])
	assert.Equal(t, nil, errs[2])
}

func TestParallelFailSafeWithLimit(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var err = errors.New("test")
	n := 100
	l := 10
	fn := 3
	trs := make([]bool, n, n)
	ts := make([]Task, n, n)

	var cc int64
	check := func() {
		assert.True(t, cc <= int64(l))
	}

	for _i := 0; _i < n; _i++ {
		i := _i
		ts[i] = TaskFunc(func(ctx context.Context) error {
			check()
			defer check()
			ta := time.After(time.Millisecond * 500)
			check()
			select {
			case <-ta:
				check()
				trs[i] = true
				check()
				return nil
			case <-ctx.Done():
				check()
				return ctx.Err()
			}
		})
	}

	ts[fn] = TaskFunc(func(ctx context.Context) error {
		trs[fn] = true
		return err
	})

	nt := ParallelFailSafeWithLimit(l, ts...)
	ne := nt.Do(context.Background())
	assert.Error(t, ne)

	for i := 0; i < n; i++ {
		assert.True(t, trs[i])
	}

	errs := ne.(Errs)
	for i := 0; i < n; i++ {
		if i == fn {
			assert.Equal(t, err, errs[i])
		} else {
			assert.Equal(t, nil, errs[i])
		}
	}
}

func TestParallelFailSafeWithLimit_EarlyCtxExit(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var done int64

	var tasks []Task
	for i := 0; i < 100; i++ {
		tasks = append(tasks, TaskFunc(func(ctx context.Context) error {
			time.Sleep(time.Second)
			atomic.AddInt64(&done, 1)
			return nil
		}))
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()

	errs := ParallelFailSafeWithLimit(10, tasks...).Do(ctx).(Errs)

	nf := 0

	for _, err := range errs {
		if err == nil {
			nf++
		} else {
			assert.Equal(t, context.Canceled, err)
		}
	}

	assert.Equal(t, 10, nf)
	assert.Equal(t, int64(10), done)
}
