package models

import (
	"context"
	"log"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/crossle/hacker-news-mixin-bot/session"
)

type Subscriber struct {
	UserId    string
	CreatedAt time.Time
}

func CreateSubscriber(ctx context.Context, userId string) (*Subscriber, error) {
	subscriber, err := findSubscriberById(ctx, userId)
	if subscriber == nil {
		if _, err := bot.UuidFromString(userId); err != nil {
			return nil, err
		}
		user := &Subscriber{
			UserId:    userId,
			CreatedAt: time.Now(),
		}
		_, err = session.Database(ctx).Exec("Insert into subscribers(user_id, created_at) values($1, $2)", user.UserId, user.CreatedAt)
		if err != nil {
			return nil, err
		}
		return user, nil

	}
	return subscriber, nil
}

func RemoveSubscriber(ctx context.Context, userId string) error {
	subscriber, err := findSubscriberById(ctx, userId)
	if subscriber == nil {
		return nil
	}
	_, err = session.Database(ctx).Exec("Delete from subscribers where user_id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}

func FindSubscribers(ctx context.Context) ([]*Subscriber, error) {
	var users []*Subscriber
	rows, err := session.Database(ctx).Query("select user_id, created_at from subscribers")
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var subscriber Subscriber
		if err := rows.Scan(&subscriber.UserId, &subscriber.CreatedAt); err != nil {
			log.Println(err.Error())
			return users, err
		}
		users = append(users, &subscriber)
	}
	return users, nil
}

func findSubscriberById(ctx context.Context, userId string) (*Subscriber, error) {
	row := session.Database(ctx).QueryRow("select user_id, created_at from subscribers where user_id = $1", userId)
	var subscriber Subscriber
	if err := row.Scan(&subscriber.UserId, &subscriber.CreatedAt); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &subscriber, nil
}
