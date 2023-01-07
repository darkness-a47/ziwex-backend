package utils

import (
	"context"
	"time"
)

func GetDatabaseContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Millisecond*200)
}
