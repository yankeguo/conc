package conc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrs_Error(t *testing.T) {
	errs := Errs{errors.New("test1"), nil, errors.New("test2")}
	assert.Equal(t, "#0: test1; #2: test2", errs.Error())
	errs = nil
	assert.Equal(t, "", errs.Error())
}

func TestErrs_Sanitize(t *testing.T) {
	errs := Errs{errors.New("test1"), nil, errors.New("test2")}
	assert.Equal(t, errs, errs.Sanitize())
	errs = Errs{nil, nil, nil}
	assert.Equal(t, nil, errs.Sanitize())
	errs = nil
	assert.Equal(t, nil, errs.Sanitize())
}

func TestErrs_Unwrap(t *testing.T) {
	err1 := errors.New("test1")
	errs := Errs{err1, nil, errors.New("test2")}
	assert.Equal(t, []error{errors.New("test1"), errors.New("test2")}, errs.Unwrap())
	assert.True(t, errors.Is(errs, err1))
	errs = nil
	assert.Equal(t, []error(nil), errs.Unwrap())
}

func TestNoError(t *testing.T) {
	et := Error(errors.New("failed"))
	err := et.Do(context.Background())
	require.Error(t, err)
	net := NoError(et)
	err = net.Do(context.Background())
	require.NoError(t, err)
}
