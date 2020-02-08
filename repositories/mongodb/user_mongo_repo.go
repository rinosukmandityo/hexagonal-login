package mongo

import (
	"context"
	"gopkg.in/mgo.v2/bson"
	"time"

	"github.com/rinosukmandityo/hexagonal-login/helper"
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

func NewUserMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (repo.LoginRepository, error) {
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

func (r *userMongoRepository) GetAll(param repo.GetAllParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(param.Tablename)
	csr, e := c.Find(ctx, nil)
	if e != nil {
		return e
	}
	if e := csr.Decode(&param.Result); e != nil {
		return e
	}
	return nil
}
func (r *userMongoRepository) GetBy(param repo.GetParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(param.Tablename)
	if e := c.FindOne(ctx, param.Filter).Decode(param.Result); e != nil {
		if e == mongo.ErrNoDocuments {
			return errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
		}
		return errors.Wrap(e, "repository.User.GetById")
	}
	return nil

}
func (r *userMongoRepository) Store(param repo.StoreParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(param.Tablename)
	if _, e := c.InsertOne(ctx, param.Data); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}

	return nil

}
func (r *userMongoRepository) Update(param repo.UpdateParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(param.Tablename)
	if res, e := c.UpdateOne(ctx, param.Filter, bson.M{"$set": param.Data}, options.Update().SetUpsert(false)); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	} else {
		if res.MatchedCount == 0 && res.ModifiedCount == 0 {
			return errors.Wrap(errors.New("User Not Found"), "repository.User.Update")
		}
	}

	return nil

}
func (r *userMongoRepository) Delete(param repo.DeleteParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(param.Tablename)
	if res, e := c.DeleteOne(ctx, param.Filter); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	} else {
		if res.DeletedCount == 0 {
			return errors.Wrap(errors.New("User Not Found"), "repository.User.Delete")
		}
	}

	return nil

}
