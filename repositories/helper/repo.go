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
	switch os.Getenv("url_db") {
	case "mysql":
		url := os.Getenv("mysql_url")
		if url == "" {
			url = "root:Password.1@tcp(127.0.0.1:3306)/tes"
		}
		db := os.Getenv("mysql_db")
		if db == "" {
			db = "tes"
		}
		timeout, _ := strconv.Atoi(os.Getenv("mysql_timeout"))
		if timeout == 0 {
			timeout = 10
		}
		repo, e := mr.NewUserRepository(url, db, timeout)
		if e != nil {
			log.Fatal(e)
		}

		return repo
	default:
		url := os.Getenv("mongo_url")
		if url == "" {
			url = "mongodb://localhost:27017/local"
		}
		db := os.Getenv("mongo_db")
		if db == "" {
			db = "local"
		}
		timeout, _ := strconv.Atoi(os.Getenv("mongo_timeout"))
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
