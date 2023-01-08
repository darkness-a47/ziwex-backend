package utils

import (
	"context"
	"time"
)

func GetDatabaseContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*200)
}

func GetMinioPutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*500)
}

func GetMinioGetContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
