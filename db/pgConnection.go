package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pg *pgxpool.Pool

func PgInit() func() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var err error
	Pg, err = pgxpool.New(ctx, "postgresql://darth:@localhost:5432/ziwex")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database:\n|	%s\n", err.Error())
		os.Exit(1)
	}
	return Pg.Close
}
