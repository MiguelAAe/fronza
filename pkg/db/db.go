package db

import (
	"context"
	"database/sql"
	"portal-backend/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func Init() (*bun.DB, error) {
	switch config.Environment {
	case "prod":
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.DatabaseURL)))
		db = bun.NewDB(sqldb, pgdialect.New())
	case "test":
		pgconn := pgdriver.NewConnector(
			pgdriver.WithNetwork("tcp"),
			pgdriver.WithAddr("db:5432"),
			pgdriver.WithUser("postgres"),
			pgdriver.WithPassword("admin"),
			pgdriver.WithInsecure(true))

		_, err := pgconn.Connect(context.Background())
		if err != nil {
			return nil, err
		}

		sqldb := sql.OpenDB(pgconn)
		db = bun.NewDB(sqldb, pgdialect.New())
	default:
		panic("environment variable not set")
	}

	err := db.Ping()

	return db, err
}
