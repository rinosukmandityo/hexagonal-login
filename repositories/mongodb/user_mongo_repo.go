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

func NewUserMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (repo.UserRepository, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	res := []m.User{}
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
func (r *userMongoRepository) GetById(id string) (*m.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	user := new(m.User)
	c := r.client.Database(r.database).Collection(user.TableName())
	if e := c.FindOne(ctx, bson.M{"_id": id}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
		}
		return user, errors.Wrap(e, "repository.User.GetById")
	}
	return user, nil

}
func (r *userMongoRepository) GetByUsername(username string) (bool, *m.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	user := new(m.User)
	c := r.client.Database(r.database).Collection(user.TableName())
	if e := c.FindOne(ctx, bson.M{"Username": username}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			return false, nil, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetById")
		}
		return false, user, errors.Wrap(e, "repository.User.GetById")
	}
	return true, user, nil

}
func (r *userMongoRepository) Store(user *m.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(new(m.User).TableName())
	if _, e := c.InsertOne(ctx, user); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}

	return nil

}
func (r *userMongoRepository) Update(user *m.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(new(m.User).TableName())
	if res, e := c.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"Username": user.Username}}, options.Update().SetUpsert(false)); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	} else {
		if res.MatchedCount == 0 && res.ModifiedCount == 0 {
			return errors.Wrap(errors.New("User Not Found"), "repository.User.Update")
		}
	}

	return nil

}
func (r *userMongoRepository) Delete(user *m.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(new(m.User).TableName())
	if res, e := c.DeleteOne(ctx, bson.M{"_id": user.ID}); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	} else {
		if res.DeletedCount == 0 {
			return errors.Wrap(errors.New("User Not Found"), "repository.User.Delete")
		}
	}

	return nil

}
