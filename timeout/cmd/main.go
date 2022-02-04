package main

import (
	"context"
	"fmt"
	"time"

	"github.com/MehdiEidi/cloud-native-patterns/timeout"
)

func main() {
	ctx := context.Background()
	dctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	timeout := timeout.Timeout(Slow)
	res, err := timeout(dctx, "some input")

	fmt.Println(res, err)
}

// Sample func
func Slow(str string) (string, error) {
	return "", nil
}
