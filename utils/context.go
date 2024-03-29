package utils

import (
	"context"
	"time"
)

func GetPgContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*500)
}

func GetMinioContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*500)
}

func GetMinioGetContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func GetRedisContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*200)
}
