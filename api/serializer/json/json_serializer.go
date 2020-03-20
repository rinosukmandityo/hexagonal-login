package json

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"

	"encoding/json"
	"github.com/pkg/errors"
)

type User struct{}

func (u *User) Decode(input []byte) (*m.User, error) {
	user := new(m.User)
	if e := json.Unmarshal(input, user); e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.Decode")
	}
	return user, nil
}

func (u *User) Encode(input *m.User) ([]byte, error) {
	rawMsg, e := json.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.logic.Encode")
	}
	return rawMsg, nil
}

func (u *User) DecodeMap(input []byte) (map[string]interface{}, error) {
	user := map[string]interface{}{}
	if e := json.Unmarshal(input, &user); e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.DecodeMap")
	}
	return user, nil
}

func (u *User) EncodeMap(input map[string]interface{}) ([]byte, error) {
	rawMsg, e := json.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.logic.EncodeMap")
	}
	return rawMsg, nil
}
