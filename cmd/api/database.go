package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"

	"github.com/lokatalent/backend_go/cmd/api/util"
)

const DB_TIMEOUT = 10 * time.Second

// openDB returns a new PostgreSQL connection pool.
func openDB(config *util.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DB.DSN)
	if err != nil {
		return nil, err
	}

	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	db = sqldblogger.OpenDriver(config.DB.DSN, db.Driver(), loggerAdapter)
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
