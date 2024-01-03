package conc

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewFixedLimiterPanic(t *testing.T) {
	assert.Panics(t, func() {
		NewFixedLimiter(-1)
	})
}

func TestNewFixedLimiter(t *testing.T) {
	rand.Seed(time.Now().Unix())
	n := rand.Intn(10) + 10
	lm := NewFixedLimiter(n)

	var c int64

	check := func() {
		assert.True(t, c <= int64(n))
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			check()
			lm.Take()
			defer lm.Done()
			check()
			atomic.AddInt64(&c, 1)
			defer atomic.AddInt64(&c, -1)
			check()
			time.Sleep(time.Millisecond * 10 * time.Duration(rand.Intn(10)+10))
			check()
		}()
	}

	wg.Wait()
}
