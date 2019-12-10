package lib

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/garyburd/redigo/redis"
	"time"
)

// redis连接配置
type RedisConfig struct {
	Host      string
	Port      int
	Password  string
	URL       string
	MaxIdle   int
	MaxActive int
}

var	RedisPool *redis.Pool

func RedisConnect(rds RedisConfig) {
	RedisPool = &redis.Pool{
		MaxIdle:     rds.MaxIdle,
		MaxActive:   rds.MaxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rds.URL, redis.DialPassword(rds.Password))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

/**
关闭连接池子
 */
func RedisClose() {
	if RedisPool != nil {
		//LogDebug("redis was closed", "application")
		_ = RedisPool.Close()
	}
}

type rds struct {
	rc redis.Conn
}


func NewRedisClient() *rds {
	return &rds{
		rc: RedisPool.Get(),
	}
}

/**
关闭连接
 */
func (r *rds) Close() {
	_ = r.rc.Close()
}


func (r *rds) Set(key string, value interface{}) error {
	key = setMd5(key)
	if _, err := r.rc.Do("Set", key, value); err != nil {
		return err
	}
	return nil
}

func (r *rds) SetExpire(key string, value interface{}, expire int) error {
	key = setMd5(key)
	if _, err := r.rc.Do("Set", key, value, "EX", expire); err != nil {
		return err
	}
	return nil
}

func (r *rds) Get(key string) (interface{}, error) {
	key = setMd5(key)
	rep, err := redis.String(r.rc.Do("Get", key))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err == nil {
		return rep, nil
	}
	return "", err
}

func (r *rds) Delete(key string) error {
	key = setMd5(key)
	if _, err := r.rc.Do("DEL", key); err != nil {
		return err
	}
	return nil
}

func setMd5(str string) string{
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}



