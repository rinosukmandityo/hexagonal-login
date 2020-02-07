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

type userRedisRepository struct {
	client *redis.Client
}

func newUserRedisClient(redisURL string) (*redis.Client, error) {
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

func NewUserRedisRepository(redisURL string) (repo.UserRepository, error) {
	repo := &userRedisRepository{}
	client, e := newUserRedisClient(redisURL)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewUserRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *userRedisRepository) generateKey(code string) string {
	return fmt.Sprintf("login<>%s", code)
}

func (r *userRedisRepository) generateUsernameKey(code string) string {
	return fmt.Sprintf("login<>username<>%s", code)
}

func (r *userRedisRepository) GetAll() ([]m.User, error) {
	res := []m.User{}
	return res, nil
}

func (r *userRedisRepository) GetById(id string) (*m.User, error) {
	user := new(m.User)
	key := r.generateKey(id)
	data, e := r.client.HGetAll(key).Result()
	if e != nil {
		return user, errors.Wrap(e, "repository.Redis.GetById")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
	}
	user.ID = data["ID"]
	user.Username = data["Username"]
	user.Email = data["Email"]
	user.Password = data["Password"]
	user.Name = data["Name"]
	user.Address = data["Address"]
	user.IsActive, _ = strconv.ParseBool(data["IsActive"])
	return user, nil
}
func (r *userRedisRepository) GetByUsername(username string) (bool, *m.User, error) {
	user := new(m.User)
	key := r.generateUsernameKey(username)
	data, e := r.client.HGetAll(key).Result()
	if e != nil {
		return false, user, errors.Wrap(e, "repository.Redis.GetById")
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
	return true, user, nil
}
func (r *userRedisRepository) Store(user *m.User) error {
	key := r.generateKey(user.ID)
	data := map[string]interface{}{
		"ID":       user.ID,
		"Username": user.Username,
		"Email":    user.Email,
		"Password": user.Password,
		"Name":     user.Name,
		"Address":  user.Address,
		"IsActive": user.IsActive,
	}
	if _, e := r.client.HMSet(key, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	keyUsername := r.generateUsernameKey((user.Username))
	if _, e := r.client.HMSet(keyUsername, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	return nil

}
func (r *userRedisRepository) Update(user *m.User) error {
	key := r.generateKey(user.ID)
	data := map[string]interface{}{
		"ID":       user.ID,
		"Username": user.Username,
		"Email":    user.Email,
		"Password": user.Password,
		"Name":     user.Name,
		"Address":  user.Address,
		"IsActive": user.IsActive,
	}
	if _, e := r.client.HMSet(key, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	keyUsername := r.generateUsernameKey((user.Username))
	if _, e := r.client.HMSet(keyUsername, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	return nil

}
func (r *userRedisRepository) Delete(user *m.User) error {
	key := r.generateKey(user.ID)
	if _, e := r.client.HDel(key).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}
	keyUsername := r.generateUsernameKey(user.Username)
	if _, e := r.client.HDel(keyUsername).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}

	return nil

}
