package redis

import (
	"github.com/garyburd/redigo/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
)

var (
	Pool *redis.Pool
)

type RedisInfo struct {
	Hostname string
	Port     int
}

func Connect(ri RedisInfo) {
	redisHost := fmt.Sprintf("%s:%d", ri.Hostname, ri.Port)
	Pool = newPool(redisHost)
	cleanupHook()
}

//create pool connection to Redis server
func newPool(server string) *redis.Pool {

	return &redis.Pool{
		MaxActive:   500,
		MaxIdle:     500,
		IdleTimeout: 5 * time.Second,

		Dial: func() (redis.Conn, error) {
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

func cleanupHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}
