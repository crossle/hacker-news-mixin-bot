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

func getTopStoryId() (int, error) {
	topStories, err := h.GetStories("top")
	if err != nil {
		return 18199708, nil
	}
	return topStories[0], nil
}

func getTopTenStories() ([]int, error) {
	stories, err := h.GetStories("top")
	if err != nil {
		return nil, err
	}
	topTen := stories[:10]
	sort.Ints(topTen)
	return topTen, nil
}

func sendTopStoryToChannel(ctx context.Context, stats *Stats) {
	prevStoryId := stats.getPrevTopStoryId()
	topTenStories, err := getTopTenStories()
	if err != nil {
		log.Printf("get Top stories")
	}
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
	storyId, err := getTopStoryId()
	if err != nil {
		return nil
	}
	stats := &Stats{storyId}
	gocron.Every(5).Minutes().Do(sendTopStoryToChannel, ctx, stats)
	<-gocron.Start()
	return nil
}
