package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func NewPostgres(connString string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), connString)
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
