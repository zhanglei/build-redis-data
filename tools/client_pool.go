package tools

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type ClientPool map[string]*redis.Pool

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func()(redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (pool ClientPool) GetServerPool(addr string) *redis.Pool {
	p, ok := pool[addr]
	if ok {
		return p
	} else {
		p = newPool(addr)
		pool[addr] = p
		return p
	}
}
