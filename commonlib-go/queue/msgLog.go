package queue

import (
	"time"
	"aimarketingCRM/commonlib-go/cache"
)

type MsgLog struct {
	cache.Cache
}

const Prefix = "queue:"

func NewMsgLog(config cache.RedisConfig) *MsgLog {
	return &MsgLog{
		cache.NewRedis(config),
	}
}

func (m *MsgLog) Get(msgId string) *MsgLog{
	if m, err := m.Cache.Get(getMsgId(msgId)); err != nil {
		if l, ok := m.(MsgLog); ok {
			return &l
		}

		return nil
	}
	return nil
}

func (m *MsgLog) Set(msgId string, expire uint8) bool {
	if _, err := m.Cache.Set(getMsgId(msgId), 1, time.Duration(expire)); err != nil {
		panic(err)
		return false
	}

	return true
}

func getMsgId(msgId string) string {
	return Prefix + msgId
}

