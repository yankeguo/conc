package conc

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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

func TestNoError(t *testing.T) {
	et := Error(errors.New("failed"))
	err := et.Do(context.Background())
	require.Error(t, err)
	net := NoError(et)
	err = net.Do(context.Background())
	require.NoError(t, err)
}
