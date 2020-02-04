package api

import (
	slz "github.com/rinosukmandityo/hexagonal-login/serializer"
	js "github.com/rinosukmandityo/hexagonal-login/serializer/json"
	ms "github.com/rinosukmandityo/hexagonal-login/serializer/msgpack"
)

var (
	ContentTypeJson    = "application/json"
	ContentTypeMsgPack = "application/x-msgpack"
)

func GetSerializer(contentType string) slz.UserSerializer {
	if contentType == ContentTypeMsgPack {
		return &ms.User{}
	}
	return &js.User{}
}
