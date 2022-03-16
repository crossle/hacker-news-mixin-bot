package main

import (
	"context"
	"encoding/base64"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/crossle/hacker-news-mixin-bot/config"
	"github.com/crossle/hacker-news-mixin-bot/models"
)

type ResponseMessage struct {
	client *bot.BlazeClient
}

func (r *ResponseMessage) OnAckReceipt(ctx context.Context, msg bot.MessageView, userID string) error {
	return nil
}
func (r *ResponseMessage) SyncAck() bool {
	return true
}

func (r *ResponseMessage) OnMessage(ctx context.Context, msg bot.MessageView, uid string) error {
	if msg.Category != bot.MessageCategorySystemAccountSnapshot && msg.Category != bot.MessageCategorySystemConversation && msg.ConversationId == bot.UniqueConversationId(config.MixinClientId, msg.UserId) {
		if msg.Category == "PLAIN_TEXT" {
			data, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				return bot.BlazeServerError(ctx, err)
			}
			if "/start" == string(data) {
				_, err = models.CreateSubscriber(ctx, msg.UserId)
				if err == nil {
					if err := r.client.SendPlainText(ctx, msg, "订阅成功"); err != nil {
						return bot.BlazeServerError(ctx, err)
					}
				}
			} else if "/stop" == string(data) {
				err = models.RemoveSubscriber(ctx, msg.UserId)
				if err == nil {
					if err := r.client.SendPlainText(ctx, msg, "已取消订阅"); err != nil {
						return bot.BlazeServerError(ctx, err)
					}
				}
			} else {
				content := `请更新 Mixin Messenger 到最新版本 0.14.1+`
				if err := r.client.SendAppButton(ctx, msg.ConversationId, msg.UserId, "点我订阅", "input:/start", "#AA4848"); err != nil {
					return bot.BlazeServerError(ctx, err)
				}
				if err := r.client.SendAppButton(ctx, msg.ConversationId, msg.UserId, "点我取消订阅", "input:/stop", "#AA4848"); err != nil {
					return bot.BlazeServerError(ctx, err)
				}
				if err = r.client.SendPlainText(ctx, msg, content); err != nil {
					return bot.BlazeServerError(ctx, err)
				}
			}
		} else {
			content := `请更新 Mixin Messenger 到最新版本 0.14.1+`
			if err := r.client.SendAppButton(ctx, msg.ConversationId, msg.UserId, "点我订阅", "input:/start", "#AA4848"); err != nil {
				return bot.BlazeServerError(ctx, err)
			}
			if err := r.client.SendAppButton(ctx, msg.ConversationId, msg.UserId, "点我取消订阅", "input:/stop", "#AA4848"); err != nil {
				return bot.BlazeServerError(ctx, err)
			}
			if err := r.client.SendPlainText(ctx, msg, content); err != nil {
				return bot.BlazeServerError(ctx, err)
			}
		}
	}
	return nil
}
