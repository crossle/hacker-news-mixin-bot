package main

import (
	"context"
	"flag"
	"log"

	"github.com/crossle/hacker-news-mixin-bot/durable"
	"github.com/crossle/hacker-news-mixin-bot/services"
)

func main() {
	service := flag.String("service", "blaze", "run a service")
	flag.Parse()
	db := durable.OpenDatabaseClient(context.Background())
	defer db.Close()

	switch *service {
	case "blaze":
		err := StartBlaze(db)
		if err != nil {
			log.Println(err)
		}
	default:
		hub := services.NewHub(db)
		err := hub.StartService(*service)
		if err != nil {
			log.Println(err)
		}
	}
}
