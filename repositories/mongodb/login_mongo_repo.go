package mongo

import (
	"context"
	"time"

	"github.com/rinosukmandityo/hexagonal-login/helper"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type loginMongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newLoginMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
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

func NewLoginMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (repo.LoginRepository, error) {
	repo := &loginMongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, e := newLoginMongoClient(mongoURL, mongoTimeout)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewLoginMongoRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *loginMongoRepository) Authenticate(username, password string) (bool, *m.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	user := new(m.User)
	c := r.client.Database(r.database).Collection(user.TableName())
	if e := c.FindOne(ctx, bson.M{"$or": []bson.M{
		{"Username": username},
		{"Email": username},
	}}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			return false, nil, errors.Wrap(helper.ErrUserNotFound, "repository.Login.Authenticate")
		}
		return false, user, errors.Wrap(e, "repository.Login.Authenticate")
	}
	if !repo.IsPasswordMatch(password, user.Password) {
		return false, user, errors.Wrap(errors.New("Password is incorrect"), "repository.Login.Authenticate")
	}

	return true, user, nil
}
