package auth

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"time"
)

type RedisSessionStore struct {
	conn *redis.Client
	ctx  context.Context
}

func CreateRedisSessionStore(conn *redis.Client) *RedisSessionStore {
	rss := RedisSessionStore{
		conn: conn,
		ctx:  context.Background(),
	}
	return &rss
}

func (rss *RedisSessionStore) AddUser(u *User) string {
	cookieString := makeSessionCookieString()

	b, err := json.Marshal(u)
	if err != nil {
		log.Err(err).Msg("Error marshalling user")
	}

	rss.conn.Set(rss.ctx, "session:"+cookieString, string(b), time.Hour*6)

	return cookieString
}

func (rss *RedisSessionStore) RemoveUser(cookie string) {
	rss.conn.Del(rss.ctx, "session:"+cookie)
}

func (rss *RedisSessionStore) GetUser(cookie string) *User {
	data, err := rss.conn.Get(rss.ctx, "session:"+cookie).Result()
	if err != nil {
		log.Err(err).Msg("Error getting user from redis")
		return nil
	}

	var u User
	err = json.Unmarshal([]byte(data), &u)
	if err != nil {
		log.Err(err).Msg("Error unmarshalling user from json")
		return nil
	}

	return &u
}
