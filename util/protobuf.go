package util

import (
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"
)

func MarshalAny(typeUrl string, message proto.Message) *any.Any {
	marshaledConfig, _ := proto.Marshal(message)

	return &any.Any{
		TypeUrl: typeUrl,
		Value:   marshaledConfig,
	}
}
