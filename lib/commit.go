package lib

import (
	"fmt"

	"github.com/deckarep/gosx-notifier"
	redis "gopkg.in/redis.v5"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

const commitsRedisKey = "watcher/commits"

type Commit struct {
	Hash    string
	Author  string
	Subject string
}

func (c Commit) Save() error {
	err := redisClient.SAdd(commitsRedisKey, c.Hash).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c Commit) SendNotification() bool {
	result, err := redisClient.SIsMember(commitsRedisKey, c.Hash).Result()
	if result || err != nil {
		return false
	}

	fmt.Printf("\r\nPushed the commit %s by %s \n", c.Hash, c.Author)
	note := gosxnotifier.NewNotification(c.Subject)
	note.Title = "The Watcher"
	note.Subtitle = c.Author
	note.Sound = gosxnotifier.Basso
	note.Sender = "com.apple.Safari"
	note.Link = "https://github.com/fiverr/5rr_v2/commit/" + c.Hash
	note.AppIcon = "assets/code.png"
	err = note.Push()
	if err != nil {
		return false
	}

	return true
}
