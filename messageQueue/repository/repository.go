package repository

import (
	"messageQueue/models"

	"github.com/gomodule/redigo/redis"
)

type RedisConfig struct {
	MaxIdle   int
	MaxActive int
	Timeout   int
	Wait      bool

	URL      string
	Port     string
	Username string
	Password string
}

type database struct {
	rdb redis.Conn
}

type QueueRepository interface {
	Delete(key string) error
	PushMessage(key string, message models.MessageIn) error
	GetMessagesFromKeys(key string) ([]string, error)
}

func NewRedisDatabase(conf *RedisConfig) (redis.Conn, error) {
	rdb := redisPool(conf).Get()
	_, err := redis.String(rdb.Do("PING"))
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func redisPool(conf *RedisConfig) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", conf.URL+":"+conf.Port)
			if err != nil {
				panic(err.Error())
			}
			return conn, err
		},
	}
}

func NewQueueRepository(rdb redis.Conn) QueueRepository {
	return &database{rdb}
}

func (c *database) CreateQueue(key string) error {
	_, err := c.rdb.Do("SET", key)
	return err
}

func (c *database) Delete(key string) error {
	_, err := c.rdb.Do("DEL", key)
	return err
}

func (c *database) PushMessage(key string, message models.MessageIn) error {
	msg := message.Desc
	_, errr := c.rdb.Do("RPUSH", key, msg)
	return errr

}

func (c *database) GetMessagesFromKeys(key string) ([]string, error) {

	value, err := redis.Strings(c.rdb.Do("LRANGE", key, -100, 100))
	if err != nil {
		return value, err
	}

	return value, nil
}
