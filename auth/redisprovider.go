package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type redisUser struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	PasswordHash string `json:"passwordhash"`

	UserGroups []string `json:"groups"`
	Type       int      `json:"type"`
}

type RedisProvider struct {
	conn *redis.Client
	ctx  context.Context
}

func CreateRedisProvider(conn *redis.Client) *RedisProvider {
	return &RedisProvider{
		conn: conn,
		ctx:  context.Background(),
	}
}

func (rp *RedisProvider) GetName() string {
	return "redis"
}

func (rp *RedisProvider) GetWritable() bool {
	return true
}

func (rp *RedisProvider) storeUser(u redisUser) error {
	data, err := json.Marshal(&u)
	if err != nil {
		log.Err(err).Str("username", u.Username).Msg("Error adding user")
		return err
	}

	err = rp.conn.Set(rp.ctx, "user:"+u.Username, string(data), 0).Err()
	if err != nil {
		log.Err(err).Str("username", u.Username).Msg("Error adding user")
	}
	return err
}

func (rp *RedisProvider) AddUser(username string, password string, usertype int) error {
	hash := hashPassword(password, getDefaultHasherConfig())
	u := redisUser{
		Username:     username,
		PasswordHash: hash,
		Type:         usertype,
	}

	return rp.storeUser(u)
}

func (rp *RedisProvider) getRedisUser(username string) (redisUser, error) {
	res, err := rp.conn.Get(rp.ctx, "user:"+username).Result()
	if err != nil {
		log.Err(err).Str("username", username).Msg("Error getting user from Redis")
		return redisUser{}, err
	}

	u := redisUser{}

	err = json.Unmarshal([]byte(res), &u)
	if err != nil {
		log.Err(err).Str("username", username).Msg("Error unmarshalling user")
		return redisUser{}, err
	}

	return u, nil
}

func (rp *RedisProvider) LoginUser(username string, password string) (User, error) {
	u, err := rp.getRedisUser(username)
	if err != nil {
		return User{}, errors.New("Error getting user from Redis")
	}

	if checkPassword(password, u.PasswordHash) {
		return User{
			Name:     u.Name,
			UserName: u.Username,
			UserType: u.Type,
			Groups:   u.UserGroups,
		}, nil
	}

	return User{}, errors.New("Unable to find user with those credentials")
}

func (rp *RedisProvider) GetUser(username string) (User, error) {
	u, err := rp.getRedisUser(username)
	if err != nil {
		return User{}, err
	}

	return User{
		Name:     u.Name,
		UserName: u.Username,
		UserType: u.Type,
		Groups:   u.UserGroups,
	}, nil
}

func (rp *RedisProvider) AddUserToGroup(username string, group string) error {
	u, err := rp.getRedisUser(username)
	if err != nil {
		return err
	}

	u.UserGroups = append(u.UserGroups, group)

	return rp.storeUser(u)
}
