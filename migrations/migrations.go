package migrations

import (
	"fmt"
	"portal-backend/config"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

func Init() (*pg.DB, error) {
	var opts *pg.Options
	var err error

	switch config.Environment {
	case "prod":
		opts, err = pg.ParseURL(config.DatabaseURL)
		if err != nil {
			return nil, err
		}
	case "test":
		opts = &pg.Options{
			Addr:     "db:5432",
			User:     "postgres",
			Password: "admin",
		}
	default:
		panic("environment variable not set")
	}

	// connect to DB
	dbConn := pg.Connect(opts)

	// return the db connection
	return dbConn, nil
}

func Run(db *pg.DB) error {
	collection := migrations.NewCollection()
	err := collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return err
	}

	_, _, err = collection.Run(db, "init")
	if err != nil {
		return err
	}

	oldVersion, newVersion, err := collection.Run(db, "up")
	if err != nil {
		return err
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("db version is %d\n", oldVersion)
	}
	return nil
}
