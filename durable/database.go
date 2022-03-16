package durable

import (
	"context"
	"database/sql"
	"log"

	"github.com/crossle/hacker-news-mixin-bot/config"
	_ "modernc.org/sqlite"
)

func OpenDatabaseClient(ctx context.Context) *sql.DB {
	db, err := sql.Open("sqlite", config.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
