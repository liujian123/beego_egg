package models

import (
	"github.com/garyburd/redigo/redis"
)


// SetKey set key value into redis
func SetKey(key string, value interface{}) {
	conn := RedisPool.Get()
	defer conn.Close()
	conn.Do("SET", key, value)
}

// GetKey get value from redis by key
func GetKey(key, field string) (r string, err error) {
	conn := RedisPool.Get() //从连接池，取一个链接
	defer conn.Close()      //函数运行结束 ，把连接放回连接池
	r, err = redis.String(conn.Do("HGet", key, field))
	if err != nil {
		return
	}
	return
}
