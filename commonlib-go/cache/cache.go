package cache

import "time"

type Cache interface{
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expire time.Duration) (bool, error)
}