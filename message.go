package main

import (
	"context"
	"encoding/base64"

	bot "github.com/MixinNetwork/bot-api-go-client"

	"github.com/crossle/hacker-news-mixin-bot/config"
	"github.com/crossle/hacker-news-mixin-bot/models"
)

type ResponseMessage struct {
}

func (r ResponseMessage) OnMessage(ctx context.Context, mc *bot.MessageContext, msg bot.MessageView, uid string) error {
	if msg.Category != bot.MessageCategorySystemAccountSnapshot && msg.Category != bot.MessageCategorySystemConversation && msg.ConversationId == bot.UniqueConversationId(config.MixinClientId, msg.UserId) {
		if msg.Category == "PLAIN_TEXT" {
			data, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				return bot.BlazeServerError(ctx, err)
			}
			if "/start" == string(data) {
				_, err = models.CreateSubscriber(ctx, msg.UserId)
				if err == nil {
					if err := bot.SendPlainText(ctx, mc, msg, "订阅成功"); err != nil {
						return bot.BlazeServerError(ctx, err)
					}
				}
			} else if "/stop" == string(data) {
				err = models.RemoveSubscriber(ctx, msg.UserId)
				if err == nil {
					if err := bot.SendPlainText(ctx, mc, msg, "已取消订阅"); err != nil {
						return bot.BlazeServerError(ctx, err)
					}
				}
			} else {
				content := `发送 /start 订阅消息
发送 /stop 取消订阅`
				if err := bot.SendPlainText(ctx, mc, msg, content); err != nil {
					return bot.BlazeServerError(ctx, err)
				}
			}
		} else {
			content := `发送 /start 订阅消息
发送 /stop 取消订阅`
			if err := bot.SendPlainText(ctx, mc, msg, content); err != nil {
				return bot.BlazeServerError(ctx, err)
			}
		}
	}
	return nil
}
