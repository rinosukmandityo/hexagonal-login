package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/rinosukmandityo/hexagonal-login/helper"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type newsMySQLRepository struct {
	url     string
	timeout time.Duration
}

func newUserClient(URL string) (*sql.DB, error) {
	db, e := sql.Open("mysql", URL)
	if e != nil {
		return nil, e
	}
	if e = db.Ping(); e != nil {
		return nil, e
	}

	return db, e
}

func (r *newsMySQLRepository) createNewTable() error {
	schema := `CREATE TABLE ` + new(m.User).TableName() + ` (
		ID VARCHAR(30) NOT NULL UNIQUE,
		Name VARCHAR(30),
		Username VARCHAR(30) NOT NULL,
		Email VARCHAR(50),
		Password VARCHAR(50),
		Address VARCHAR(50),
		IsActive boolean
	);`
	db, e := sqlx.Connect("mysql", r.url)
	if e != nil {
		return errors.Wrap(e, "repository.User.CreateTable")
	}
	defer db.Close()
	res, e := db.Exec(schema)
	if res != nil && e == nil {
		fmt.Println("Table 'Users' created")
	}
	return nil
}

func NewUserRepository(URL, DB string, timeout int) (repo.UserRepository, error) {
	repo := &newsMySQLRepository{
		url:     fmt.Sprintf("%s?parseTime=true", URL),
		timeout: time.Duration(timeout) * time.Second,
	}
	repo.createNewTable()
	return repo, nil
}

func (r *newsMySQLRepository) GetAll() ([]m.User, error) {
	res := []m.User{}
	db, e := sqlx.Connect("mysql", r.url)
	if e != nil {
		return res, errors.Wrap(e, "repository.User.GetAll")
	}
	defer db.Close()
	q := constructGetAll()

	if e = db.Select(&res, q); e != nil {
		return res, errors.Wrap(e, "repository.User.GetAll")
	}
	return res, nil
}
func (r *newsMySQLRepository) GetBy(filter map[string]interface{}) (*m.User, error) {
	res := new(m.User)
	db, e := sqlx.Connect("mysql", r.url)
	if e != nil {
		return res, errors.Wrap(e, "repository.User.GetBy")
	}
	defer db.Close()
	q, dataFields := constructGetBy(filter)

	if e = db.Get(res, q, dataFields...); e != nil {
		if e == sql.ErrNoRows {
			return res, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetBy")
		}
		return res, errors.Wrap(e, "repository.User.GetBy")
	}
	return res, nil

}
func (r *newsMySQLRepository) Store(data *m.User) error {
	db, e := newUserClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}

	q, dataField := constructStoreQuery(data)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	defer stmt.Close()
	if _, e := stmt.Exec(dataField...); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}

	return nil

}

func (r *newsMySQLRepository) Update(data *m.User) error {
	db, e := newUserClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	defer conn.Close()

	filter := map[string]interface{}{"ID": data.ID}
	q, dataField := constructUpdateQuery(data, filter)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.User.Update")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(dataField...); e != nil {
		return errors.Wrap(e, "repository.User.Update")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return errors.Wrap(e, "repository.User.Update")
		}
		if count == 0 {
			return errors.Wrap(helper.ErrUserNotFound, "repository.User.Update")
		}
	}

	return nil

}
func (r *newsMySQLRepository) Delete(data *m.User) error {
	db, e := newUserClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}
	defer conn.Close()

	filter := map[string]interface{}{"ID": data.ID}
	q, dataFields := constructDeleteQuery(filter)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(dataFields...); e != nil {
		return errors.Wrap(e, "repository.User.Delete")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return errors.Wrap(e, "repository.User.Delete")
		}
		if count == 0 {
			return errors.Wrap(helper.ErrUserNotFound, "repository.User.Delete")
		}
	}

	return nil

}

func (r *newsMySQLRepository) Authenticate(username, password string) (bool, *m.User, error) {
	res := new(m.User)
	db, e := sqlx.Connect("mysql", r.url)
	if e != nil {
		return false, res, errors.Wrap(e, "repository.User.Authenticate")
	}
	defer db.Close()
	q, dataFields := constructAuth(map[string]interface{}{"Username": username, "Email": username})

	if e = db.Get(res, q, dataFields...); e != nil {
		return false, res, errors.Wrap(e, "repository.User.Authenticate")
	}
	if res.ID == "" {
		return false, res, errors.Wrap(helper.ErrUserNotFound, "repository.User.Authenticate")
	}

	if !repo.IsPasswordMatch(password, res.Password) {
		return false, res, errors.New("Password does not match")
	}

	return true, res, nil
}
