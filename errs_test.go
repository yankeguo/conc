package conc

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
