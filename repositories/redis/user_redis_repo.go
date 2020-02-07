package redis

import (
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

func (r *userRedisRepository) GetAll() ([]m.User, error) {
	res := []m.User{}
	return res, nil
}

func (r *userRedisRepository) GetById(id string) (*m.User, error) {
	user := new(m.User)
	key := generateKey(id)
	data, e := r.client.HGetAll(key).Result()
	if e != nil {
		return user, errors.Wrap(e, "repository.User.GetById")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
	}
	user.FormingUserData(data)
	return user, nil
}
func (r *userRedisRepository) GetByUsername(username string) (bool, *m.User, error) {
	user := new(m.User)
	key := generateUsernameKey(username)
	data, e := r.client.HGetAll(key).Result()
	if e != nil {
		return false, user, errors.Wrap(e, "repository.User.GetById")
	}
	if len(data) == 0 {
		return false, nil, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
	}
	user.FormingUserData(data)
	return true, user, nil
}
func (r *userRedisRepository) Store(user *m.User) error {
	key := generateKey(user.ID)
	user.Password = repo.EncryptPassword(user.Password)
	data := user.FormingData()
	if _, e := r.client.HMSet(key, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	keyUsername := generateUsernameKey((user.Username))
	if _, e := r.client.HMSet(keyUsername, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	return nil

}
func (r *userRedisRepository) Update(user *m.User) error {
	key := generateKey(user.ID)
	data := user.FormingData()
	if _, e := r.client.HMSet(key, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	keyUsername := generateUsernameKey((user.Username))
	if _, e := r.client.HMSet(keyUsername, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	return nil

}
func (r *userRedisRepository) Delete(user *m.User) error {
	key := generateKey(user.ID)
	if _, e := r.client.HDel(key).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}
	keyUsername := generateUsernameKey(user.Username)
	if _, e := r.client.HDel(keyUsername).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}

	return nil

}
