package helper

import (
	"log"
	"os"
	"strconv"

	repo "github.com/rinosukmandityo/hexagonal-login/repositories"
	mr "github.com/rinosukmandityo/hexagonal-login/repositories/mongodb"
	rr "github.com/rinosukmandityo/hexagonal-login/repositories/redis"
)

func ChooseRepo() repo.UserRepository {
	switch os.Getenv("url_db") {
	case "redis":
		redisURL := os.Getenv("redis_url")
		repo, e := rr.NewUserRedisRepository(redisURL)
		if e != nil {
			log.Fatal(e)
		}
		return repo
	default:
		mongoURL := os.Getenv("mongo_url")
		if mongoURL == "" {
			mongoURL = "mongodb://localhost:27017/local"
		}
		mongoDB := os.Getenv("mongo_db")
		if mongoDB == "" {
			mongoDB = "local"
		}
		mongoTimeout, _ := strconv.Atoi(os.Getenv("mongo_timeout"))
		if mongoTimeout == 0 {
			mongoTimeout = 10
		}

		repo, e := mr.NewUserMongoRepository(mongoURL, mongoDB, mongoTimeout)
		if e != nil {
			log.Fatal(e)
		}
		return repo
	}
	return nil
}
