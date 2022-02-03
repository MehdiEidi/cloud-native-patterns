package cb

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int
	lastAttempt := time.Now()
	var lock sync.RWMutex

	return func(ctx context.Context) (string, error) {
		lock.RLock()

		d := consecutiveFailures - int(failureThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				lock.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		lock.RUnlock()

		res, err := circuit(ctx)

		lock.Lock()
		defer lock.Unlock()

		lastAttempt = time.Now()

		if err != nil {
			consecutiveFailures++
			return res, err
		}

		consecutiveFailures = 0

		return res, nil
	}
}
