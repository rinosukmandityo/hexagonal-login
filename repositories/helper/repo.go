package repohelper

import (
	"log"
	"os"
	"strconv"

	repo "github.com/rinosukmandityo/hexagonal-login/repositories"
	mg "github.com/rinosukmandityo/hexagonal-login/repositories/mongodb"
	mr "github.com/rinosukmandityo/hexagonal-login/repositories/mysql"
)

func ChooseRepo() repo.UserRepository {
	url := os.Getenv("url")
	db := os.Getenv("db")
	timeout, _ := strconv.Atoi(os.Getenv("timeout"))
	switch os.Getenv("driver") {
	case "mysql":
		if url == "" {
			url = "root:Password.1@tcp(127.0.0.1:3306)/tes"
		}
		if db == "" {
			db = "tes"
		}
		if timeout == 0 {
			timeout = 10
		}
		repo, e := mr.NewUserRepository(url, db, timeout)
		if e != nil {
			log.Fatal(e)
		}

		return repo
	default:
		if url == "" {
			url = "mongodb://localhost:27017/local"
		}
		if db == "" {
			db = "local"
		}
		if timeout == 0 {
			timeout = 10
		}

		repo, e := mg.NewUserMongoRepository(url, db, timeout)
		if e != nil {
			log.Fatal(e)
		}
		return repo
	}
	return nil
}
