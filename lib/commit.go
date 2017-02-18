package lib

import (
	"fmt"
	"os/exec"

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
		exec.Command("osx-notifier", "--type", "info", "--message", c.Subject, "--title", "The Watcher", "--open", "https://github.com/fiverr/5rr_v2/commit/"+c.Hash, "--subtitle", c.Author).Run()
		return true
	}
}
