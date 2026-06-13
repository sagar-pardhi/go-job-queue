package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(connString string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), connString)
}

func BuildConnectionString(
	host,
	port,
	user,
	password,
	dbname string,
) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)
}
