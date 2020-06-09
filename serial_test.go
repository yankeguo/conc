package conc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSerial(t *testing.T) {
	var err = errors.New("test")
	var t1r, t2r, t3r bool
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 10)
		t1r = true
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 20)
		t2r = true
		return err
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 20)
		t3r = true
		return err
	}
	nt := Serial(t1, t2, t3)
	ne := nt.Do(context.Background())
	assert.Error(t, ne)
	assert.True(t, t1r)
	assert.True(t, t2r)
	assert.False(t, t3r)

	errs := ne.(Errs)
	assert.Equal(t, nil, errs[0])
	assert.Equal(t, err, errs[1])
	assert.Equal(t, nil, errs[2])
}

func TestSerialFailSafe(t *testing.T) {
	var err = errors.New("test")
	var t1r, t2r, t3r bool
	var t1 TaskFunc = func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 10)
		t1r = true
		return nil
	}
	var t2 TaskFunc = func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 20)
		t2r = true
		return err
	}
	var t3 TaskFunc = func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 20)
		t3r = true
		return err
	}
	nt := SerialFailSafe(t1, t2, t3)
	ne := nt.Do(context.Background())
	assert.Error(t, ne)
	assert.True(t, t1r)
	assert.True(t, t2r)
	assert.True(t, t3r)

	errs := ne.(Errs)
	assert.Equal(t, nil, errs[0])
	assert.Equal(t, err, errs[1])
	assert.Equal(t, err, errs[2])
}
