package durable

import (
	"context"
	"database/sql"
	"log"

	"github.com/crossle/hacker-news-mixin-bot/config"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabaseClient(ctx context.Context) *sql.DB {
	db, err := sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
