package debounce

import (
	"context"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	var threshold time.Time
	var result string
	var err error
	var lock sync.Mutex

	return func(ctx context.Context) (string, error) {
		lock.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			lock.Unlock()
		}()

		if time.Now().Before(threshold) {
			return result, err
		}

		result, err = circuit(ctx)

		return result, err
	}
}
