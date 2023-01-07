package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Poll *pgxpool.Pool

func ConnetionInit() func() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var err error
	Poll, err = pgxpool.New(ctx, "postgresql://darth:@localhost:5432/ziwex")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database:\n|	%s\n", err.Error())
		os.Exit(1)
	}
	return Poll.Close
}
