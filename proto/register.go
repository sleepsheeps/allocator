package proto

import (
	"github.com/golang/protobuf/proto"
	"github.com/sleepsheeps/allocator/proto/common"
	"github.com/sleepsheeps/allocator/proto/helper"
)

var (
	register = map[helper.MsgType]proto.Message{
		helper.MsgType_Transformation: &common.Transformation{},
		helper.MsgType_HeartBeat:      &common.Heartbeat{},
	}
)

func NewMsg(msgType helper.MsgType) (proto.Message, bool) {
	msg, ok := register[msgType]
	if !ok {
		return nil, false
	}
	return msg, ok
}
