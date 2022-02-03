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

func DebounceLast(circuit Circuit, d time.Duration) Circuit {
	threshold := time.Now()
	var ticker *time.Ticker
	var result string
	var err error
	var once sync.Once
	var lock sync.Mutex

	return func(ctx context.Context) (string, error) {
		lock.Lock()
		defer lock.Unlock()

		threshold = time.Now().Add(d)

		once.Do(func() {
			ticker = time.NewTicker(time.Millisecond * 100)

			go func() {
				defer func() {
					lock.Lock()
					ticker.Stop()
					once = sync.Once{}
					lock.Unlock()
				}()

				for {
					select {
					case <-ticker.C:
						lock.Lock()
						if time.Now().After(threshold) {
							result, err = circuit(ctx)
							lock.Unlock()
							return
						}
						lock.Unlock()

					case <-ctx.Done():
						lock.Lock()
						result, err = "", ctx.Err()
						lock.Unlock()
						return
					}
				}
			}()
		})

		return result, err
	}
}
