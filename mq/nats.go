package mq

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	protos "github.com/sleepsheeps/allocator/proto"
	"github.com/sleepsheeps/allocator/proto/common"
	"github.com/sleepsheeps/allocator/proto/helper"
	"log"
)

var (
	gNats *nats.Conn
	gSub  *nats.Subscription
)

func Init(url, sub string) (err error) {
	if sub == "" {
		return errors.New("nil sub")
	}
	if url == "" {
		url = nats.DefaultURL
	}
	gNats, err = nats.Connect(url)
	if err != nil {
		return err
	}
	gSub, err = gNats.Subscribe(sub, handleTransMsg)
	if err != nil {
		return err
	}
	return nil
}

func handleTransMsg(p *nats.Msg) {
	d := p.Data
	// marshal trans msg
	commonP, ok := protos.NewMsg(helper.MsgType_Transformation)
	if !ok {
		panic("not register common msg")
	}
	commonM, ok := commonP.(*common.Transformation)
	if !ok {
		log.Fatalf("transform common msg fail")
		return
	}
	if err := commonM.Unmarshal(d); err != nil {
		log.Fatalf(err.Error())
		return
	}
	if handle, ok := registerMap[commonM.GetMsgType()]; ok {
		handleP, ok := protos.NewMsg(commonM.GetMsgType())
		if !ok {
			log.Fatalf("not find method %v", commonM.GetMsgType())
			return
		}
		err := handleP.Unmarshal(commonM.GetMsgContent())
		if err != nil {
			log.Fatalf("unmarshal %v error %v", commonM.GetMsgType(), err.Error())
			return
		}
		if err = handle(handleP); err != nil {
			log.Fatalf("handle %v error %v", commonM.GetMsgType(), err.Error())
			return
		}
		log.Fatalf("handle %v ok", commonM.GetMsgType().String())
		return
	}
	log.Fatalf("not find register method %v", commonM.GetMsgType())
}

type registerFunc func(p protos.MyProto) error

var (
	registerMap = make(map[helper.MsgType]registerFunc)
)

func RegisterMsg(msgType helper.MsgType, f registerFunc) {
	_, ok := registerMap[msgType]
	if ok {
		panic(fmt.Errorf("register repeated %v", msgType))
	}
	registerMap[msgType] = f
}
