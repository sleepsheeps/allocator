package protos

import (
	"github.com/golang/protobuf/proto"
	"github.com/sleepsheeps/allocator/proto/common"
	"github.com/sleepsheeps/allocator/proto/helper"
)

var (
	register = map[helper.MsgType]MyProto{
		helper.MsgType_Transformation: &common.Transformation{},
		helper.MsgType_HeartBeat:      &common.Heartbeat{},
	}
)

type MyProto interface {
	proto.Message
	MarshalTo(dAtA []byte) (int, error)
	Unmarshal(dAtA []byte) error
}

func NewMsg(msgType helper.MsgType) (MyProto, bool) {
	msg, ok := register[msgType]
	if !ok {
		return nil, false
	}
	return msg, ok
}
