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

func (c Commit) Save() {
	redisClient.SAdd(commitsRedisKey, c.Hash)
}

func (c Commit) SendNotification() bool {
	result := redisClient.SIsMember(commitsRedisKey, c.Hash).Val()
	if result {
		return false
	} else {
		fmt.Printf("Pushed the commit %s by %s \n", c.Hash, c.Author)

		note := gosxnotifier.NewNotification(c.Subject)

		note.Title = "The Watcher"

		//Optionally, set a subtitle
		note.Subtitle = c.Author

		//Optionally, set a sound from a predefined set.
		note.Sound = gosxnotifier.Basso

		//Optionally, specifiy a url or bundleid to open should the notification be
		//clicked.
		note.Sender = "com.apple.Safari"
		note.Link = "https://github.com/fiverr/5rr_v2/commit/" + c.Hash

		//Optionally, an app icon (10.9+ ONLY)
		note.AppIcon = "assets/code.png"

		//Then, push the notification
		note.Push()

		return true
	}
}
