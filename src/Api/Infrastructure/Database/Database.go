package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitializeDatabaseConnection() {
	once.Do(func() {
		var err error

		dsn, exists := os.LookupEnv("DATABASE_DSN")

		if !exists {
			slog.Error("missing env variable: DATABASE_DSN")
			panic("missing env variable: DATABASE_DSN")
		}

		db, err = sql.Open("postgres", dsn)

		if err != nil {
			message := fmt.Sprintf("unable to open connection to database. error: %s", err.Error())
			slog.Error(message)
			panic(message)
		}

		context, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		err = db.PingContext(context)
		cancel()

		if err != nil {
			db.Close()
			message := fmt.Sprintf("unable to ping database. error: %s", err.Error())
			slog.Error(message)
			panic(message)
		}
	})
}

func GetDatabaseConnection() *sql.DB {
	InitializeDatabaseConnection()

	return db
}
