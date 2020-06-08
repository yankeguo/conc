package conc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskFunc_Run(t *testing.T) {
	err := errors.New("test")
	invoked := false

	var f TaskFunc = func(ctx context.Context) error {
		invoked = true
		return err
	}

	assert.Equal(t, err, f.Run(context.Background()))
	assert.True(t, invoked)
}
