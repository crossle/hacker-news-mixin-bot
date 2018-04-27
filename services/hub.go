package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/crossle/hacker-news-mixin-bot/durable"
	"github.com/crossle/hacker-news-mixin-bot/session"
)

type Hub struct {
	context  context.Context
	services map[string]Service
}

func NewHub(db *sql.DB) *Hub {
	hub := &Hub{services: make(map[string]Service)}
	hub.context = session.WithDatabase(context.Background(), db)
	hub.registerServices()
	return hub
}

func (hub *Hub) StartService(name string) error {
	service := hub.services[name]
	if service == nil {
		return fmt.Errorf("no service found: %s", name)
	}

	logger, err := durable.NewLoggerClient("", true)
	if err != nil {
		return err
	}
	defer logger.Close()
	ctx := session.WithLogger(hub.context, durable.BuildLogger(logger, name, nil))

	return service.Run(ctx)
}

func (hub *Hub) registerServices() {
	hub.services["news"] = &NewsService{}
}
