package mongo

import (
	"context"
	"time"

	"github.com/rinosukmandityo/hexagonal-login/helper"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type userMongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newUserMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, e := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if e != nil {
		return nil, e
	}
	if e = client.Ping(ctx, readpref.Primary()); e != nil {
		return nil, e
	}
	return client, e
}

func NewUserRepository(mongoURL, mongoDB string, mongoTimeout int) (repo.UserRepository, error) {
	repo := &userMongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, e := newUserMongoClient(mongoURL, mongoTimeout)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewUserMongoRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *userMongoRepository) GetAll() ([]m.User, error) {
	res := []m.User{}
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(new(m.User).TableName())
	csr, e := c.Find(ctx, nil)
	if e != nil {
		return res, e
	}
	if e := csr.Decode(&res); e != nil {
		return res, e
	}
	return res, nil
}
func (r *userMongoRepository) GetBy(filter map[string]interface{}) (*m.User, error) {
	res := new(m.User)
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(res.TableName())
	convertID(filter)
	if e := c.FindOne(ctx, filter).Decode(res); e != nil {
		if e == mongo.ErrNoDocuments {
			return res, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetBy")
		}
		return res, errors.Wrap(e, "repository.User.GetBy")
	}
	return res, nil

}
func (r *userMongoRepository) Store(data *m.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(data.TableName())
	if _, e := c.InsertOne(ctx, data); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}

	return nil

}
func (r *userMongoRepository) Update(data map[string]interface{}, id string) (*m.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(new(m.User).TableName())

	user := new(m.User)
	var e error

	convertID(data)
	filter := map[string]interface{}{"_id": id}

	if res, e := c.UpdateOne(ctx, filter, map[string]interface{}{"$set": data}, options.Update().SetUpsert(false)); e != nil {
		return user, errors.Wrap(e, "repository.User.Update")
	} else {
		if res.MatchedCount == 0 && res.ModifiedCount == 0 {
			return user, errors.Wrap(errors.New("User Not Found"), "repository.User.Update")
		}
	}
	user, e = r.GetBy(filter)
	if e != nil {
		return user, errors.Wrap(e, "repository.User.Update")
	}

	return user, nil
}
func (r *userMongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	filter := map[string]interface{}{"_id": id}

	c := r.client.Database(r.database).Collection(new(m.User).TableName())
	if res, e := c.DeleteOne(ctx, filter); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	} else {
		if res.DeletedCount == 0 {
			return errors.Wrap(errors.New("User Not Found"), "repository.User.Delete")
		}
	}

	return nil
}

func (r *userMongoRepository) Authenticate(username, password string) (bool, *m.User, error) {
	user := new(m.User)
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(user.TableName())
	if e := c.FindOne(ctx, map[string]interface{}{"$or": []map[string]interface{}{
		{"Username": username},
		{"Email": username},
	}}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			return false, user, errors.Wrap(helper.ErrUserNotFound, "repository.User.Authenticate")
		}
		return false, user, errors.Wrap(e, "repository.User.Authenticate")
	}
	if !repo.IsPasswordMatch(password, user.Password) {
		return false, user, errors.New("Password does not match")
	}

	return true, user, nil
}

func convertID(data map[string]interface{}) {
	if _, ok := data["ID"]; ok {
		data["_id"] = data["ID"]
		delete(data, "ID")
	}
}
