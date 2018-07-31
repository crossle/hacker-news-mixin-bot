package services

import (
	"context"
	"encoding/base64"
	"log"
	"sort"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/crossle/hacker-news-mixin-bot/config"
	"github.com/crossle/hacker-news-mixin-bot/models"
	"github.com/jasonlvhit/gocron"
	h "github.com/qube81/hackernews-api-go"
)

type NewsService struct{}

type Stats struct {
	prevStoryId int
}

func (self Stats) getPrevTopStoryId() int {
	return self.prevStoryId
}

func (self *Stats) updatePrevTopStoryId(id int) {
	self.prevStoryId = id
}

func getTopStoryId() int {
	topStories, _ := h.GetStories("top")
	return topStories[0]
}

func getTopTenStories() []int {
	stories, _ := h.GetStories("top")
	topTen := stories[:10]
	sort.Ints(topTen)
	return topTen
}

func sendTopStoryToChannel(ctx context.Context, stats *Stats) {
	prevStoryId := stats.getPrevTopStoryId()
	topTenStories := getTopTenStories()
	for _, storyId := range topTenStories {
		if storyId > prevStoryId {
			story, _ := h.GetItem(storyId)
			log.Printf("Sending top story to channel...")
			stats.updatePrevTopStoryId(story.ID)
			subscribers, _ := models.FindSubscribers(ctx)
			for _, subscriber := range subscribers {
				conversationId := bot.UniqueConversationId(config.MixinClientId, subscriber.UserId)
				data := base64.StdEncoding.EncodeToString([]byte(story.Title + " " + story.URL))
				bot.PostMessage(ctx, conversationId, subscriber.UserId, bot.UuidNewV4().String(), "PLAIN_TEXT", data, config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey)
			}
		} else {
			log.Printf("Same top story ID: %d, no message sent.", storyId)
		}
	}
}
func (service *NewsService) Run(ctx context.Context) error {
	stats := &Stats{getTopStoryId()}
	gocron.Every(5).Minutes().Do(sendTopStoryToChannel, ctx, stats)
	<-gocron.Start()
	return nil
}
