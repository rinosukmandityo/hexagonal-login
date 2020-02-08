package redis

import (
	"github.com/rinosukmandityo/hexagonal-login/helper"
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

func NewUserRedisRepository(redisURL string) (repo.LoginRepository, error) {
	repo := &userRedisRepository{}
	client, e := newUserRedisClient(redisURL)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewUserRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *userRedisRepository) GetAll(param repo.GetAllParam) error {
	return nil
}

func (r *userRedisRepository) GetBy(param repo.GetParam) error {
	key := generateKey(param.Filter["_id"].(string))
	data, e := r.client.HGetAll(key).Result()
	if e != nil {
		return errors.Wrap(e, "repository.User.GetById")
	}
	if len(data) == 0 {
		return errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
	}
	// user = FormingData(data)
	return nil
}
func (r *userRedisRepository) Store(param repo.StoreParam) error {
	data := param.Data.(map[string]interface{})
	key := generateKey(data["_id"].(string))
	if _, e := r.client.HMSet(key, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	// keyUsername := generateUsernameKey((user.Username))
	// if _, e := r.client.HMSet(keyUsername, data).Result(); e != nil {
	// 	return errors.Wrap(e, "repository.User.Store")
	// }
	return nil

}
func (r *userRedisRepository) Update(param repo.UpdateParam) error {
	data := param.Data.(map[string]interface{})
	key := generateKey(data["_id"].(string))
	if _, e := r.client.HMSet(key, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	keyUsername := generateUsernameKey(data["Username"].(string))
	if _, e := r.client.HMSet(keyUsername, data).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	return nil

}
func (r *userRedisRepository) Delete(param repo.DeleteParam) error {
	key := generateKey(param.Filter["_id"].(string))
	if _, e := r.client.HDel(key).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}
	// keyUsername := generateUsernameKey(user.Username)
	// if _, e := r.client.HDel(keyUsername).Result(); e != nil {
	// 	return errors.Wrap(e, "repository.User.Delete")
	// }

	return nil

}
