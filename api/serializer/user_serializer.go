package serializer

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
)

type UserSerializer interface {
	Decode(input []byte) (*m.User, error)
	Encode(input *m.User) ([]byte, error)
	DecodeMap(input []byte) (map[string]interface{}, error)
	EncodeMap(input map[string]interface{}) ([]byte, error)
}
