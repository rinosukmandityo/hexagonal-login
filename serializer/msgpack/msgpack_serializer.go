package msgpack

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type User struct{}

func (u *User) Decode(input []byte) (*m.User, error) {
	user := new(m.User)
	if e := msgpack.Unmarshal(input, user); e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.Decode")
	}
	return user, nil
}

func (u *User) Encode(input *m.User) ([]byte, error) {
	rawMsg, e := msgpack.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.logic.Encode")
	}
	return rawMsg, nil
}
