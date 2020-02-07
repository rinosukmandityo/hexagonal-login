package redis

import (
	"fmt"
	"strconv"

	"github.com/rinosukmandityo/hexagonal-login/helper"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type loginRedisRepository struct {
	client *redis.Client
}

func newLoginRedisClient(redisURL string) (*redis.Client, error) {
	opt, e := redis.ParseURL(redisURL)
	if e != nil {
		return nil, e
	}
	client := redis.NewClient(opt)
	if _, e = client.Ping().Result(); e != nil {
		return nil, e
	}
	return client, e
}

func NewLoginRedisRepository(redisURL string) (repo.LoginRepository, error) {
	repo := &loginRedisRepository{}
	client, e := newLoginRedisClient(redisURL)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewLoginRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *loginRedisRepository) generateUsernameKey(code string) string {
	return fmt.Sprintf("login<>username<>%s", code)
}

func (r *loginRedisRepository) Authenticate(username, password string) (bool, *m.User, error) {
	user := new(m.User)
	key := generateUsernameKey(username)
	data, e := r.client.HGetAll(key).Result()
	if e != nil {
		return false, user, errors.Wrap(e, "repository.User.GetById")
	}
	if len(data) == 0 {
		return false, nil, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
	}
	user.ID = data["ID"]
	user.Username = data["Username"]
	user.Email = data["Email"]
	user.Password = data["Password"]
	user.Name = data["Name"]
	user.Address = data["Address"]
	user.IsActive, _ = strconv.ParseBool(data["IsActive"])

	if repo.IsPasswordMatch(password, user.Password) {
		return false, user, errors.Wrap(errors.New("Password is incorrect"), "repository.Login.Authenticate")
	}
	return true, user, nil
}
