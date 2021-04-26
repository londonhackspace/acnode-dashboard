package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type redisUser struct {
	Username string `json:"username"`
	Name string `json:"name"`
	PasswordHash string `json:"passwordhash"`

	UserGroups map[string]string `json:"groups"`
	Type int `json:"type"`
}

type RedisProvider struct {
	conn *redis.Client
	ctx context.Context
}

func CreateRedisProvider(conn *redis.Client) *RedisProvider {
	return &RedisProvider{
		conn: conn,
		ctx: context.Background(),
	}
}

func (rp *RedisProvider) GetName() string {
	return "redis"
}

func (rp *RedisProvider) GetWritable() bool {
	return true
}

func (rp *RedisProvider) AddUser(username string, password string, usertype int) error {
	hash := hashPassword(password, getDefaultHasherConfig())
	u := redisUser{
		Username: username,
		PasswordHash: hash,
		Type: usertype,
	}

	data, err := json.Marshal(&u)
	if err != nil {
		log.Err(err).Str("username", username).Msg("Error adding user")
		return err
	}

	err = rp.conn.Set(rp.ctx, "user:"+username, string(data), 0).Err()
	if err != nil {
		log.Err(err).Str("username", username).Msg("Error adding user")
	}
	return err
}

func (rp *RedisProvider) LoginUser(username string, password string) (User,error) {
	res, err := rp.conn.Get(rp.ctx, "user:"+username).Result()
	if err != nil {
		log.Err(err).Str("username", username).Msg("Error getting user from Redis")
		return User{},errors.New("Error getting user from Redis")
	}

	u := redisUser{}

	err = json.Unmarshal([]byte(res), &u)
	if err != nil {
		log.Err(err).Str("username", username).Msg("Error unmarshalling user")
		return User{},errors.New("Error getting user from Redis")
	}

	if checkPassword(password, u.PasswordHash) {
		return User{
			Name: u.Name,
			UserType: u.Type,
		}, nil
	}

	return User{}, errors.New("Unable to find user with those credentials")
}