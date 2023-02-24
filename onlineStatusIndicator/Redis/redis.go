package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type RClient struct {
	client *redis.Client
}

func (rc *RClient) Init(address string) {
	rc.client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
}

func (rc *RClient) OnlineUser(usrName string) error {
	status := rc.client.Set(usrName, true, 10*time.Second)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (rc *RClient) GetOnlineUsers() ([]string, error) {
	response := rc.client.Keys("*")
	if response.Err() != nil {
		return []string{}, response.Err()
	}
	return response.Val(), nil
}
