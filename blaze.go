package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/crossle/hacker-news-mixin-bot/config"
	"github.com/crossle/hacker-news-mixin-bot/durable"
	"github.com/crossle/hacker-news-mixin-bot/session"

	bot "github.com/MixinNetwork/bot-api-go-client"
)

func StartBlaze(db *sql.DB) error {
	log.Println("start blaze")
	logger, err := durable.NewLoggerClient("", true)
	if err != nil {
		return err
	}
	defer logger.Close()
	ctx, cancel := newBlazeContext(db, logger)
	defer cancel()

	for {
		if err := bot.Loop(ctx, ResponseMessage{}, config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey); err != nil {
			session.Logger(ctx).Error(err)
		}
		session.Logger(ctx).Info("connection loop end")
		time.Sleep(time.Second)
	}
	return nil
}

func newBlazeContext(db *sql.DB, client *durable.LoggerClient) (context.Context, context.CancelFunc) {
	ctx := session.WithLogger(context.Background(), durable.BuildLogger(client, "blaze", nil))
	ctx = session.WithDatabase(ctx, db)
	return context.WithCancel(ctx)
}
