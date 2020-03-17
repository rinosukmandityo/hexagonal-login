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
	return []m.User{}, nil
}

func (r *userRedisRepository) GetBy(filter map[string]interface{}) (*m.User, error) {
	key := generateKey(filter["_id"].(string))
	data, e := r.client.HGetAll(key).Result()
	user := new(m.User)
	if e != nil {
		return user, errors.Wrap(e, "repository.User.GetById")
	}
	if len(data) == 0 {
		return user, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
	}
	user.FormingUserData(data)
	return user, nil
}
func (r *userRedisRepository) Store(data *m.User) error {
	key := generateKey(data.ID)
	redisData := data.GetMapFormat()
	if _, e := r.client.HMSet(key, redisData).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	keyUsername := generateUsernameKey(data.Username)
	if _, e := r.client.HMSet(keyUsername, redisData).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	return nil

}
func (r *userRedisRepository) Update(data *m.User, filter map[string]interface{}) error {
	key := generateKey(filter["_id"].(string))
	redisData := data.GetMapFormat()
	if _, e := r.client.HMSet(key, redisData).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	keyUsername := generateUsernameKey(filter["username"].(string))
	if _, e := r.client.HMSet(keyUsername, redisData).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	return nil

}
func (r *userRedisRepository) Delete(filter map[string]interface{}) error {
	key := generateKey(filter["_id"].(string))
	data, e := r.client.HGetAll(key).Result()
	if e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}
	if len(data) == 0 {
		return errors.Wrap(errors.New("User Not Found"), "repository.User.Delete")
	}
	if _, e := r.client.HDel(key).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}
	keyUsername := generateUsernameKey(data["username"])
	if _, e := r.client.HDel(keyUsername).Result(); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}

	return nil
}

func (r *userRedisRepository) Authenticate(username, password string) (bool, *m.User, error) {
	user := new(m.User)
	key := generateUsernameKey(username)
	data, e := r.client.HGetAll(key).Result()
	if e != nil {
		return false, user, errors.Wrap(e, "repository.User.Authenticate")
	}
	if len(data) == 0 {
		return false, user, errors.Wrap(errors.New("User Not Found"), "repository.User.Authenticate")
	}
	if !repo.IsPasswordMatch(password, user.Password) {
		return false, user, errors.New("Password does not match")
	}
	user.FormingUserData(data)

	return true, user, nil
}
