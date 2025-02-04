package main

import (
	"errors"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	return &RedisStore{Client: client}
}

func (app *application) SaveURL(originURL string) (string, error) {
	shortURL, err := app.storage.Client.Get(app.ctx, originURL).Result()
	if errors.Is(err, redis.Nil) {
		shortURL = app.createShortURL(originURL)
		err = app.storage.Client.Set(app.ctx, originURL, shortURL, 0).Err()
		if err != nil {
			return "", err
		}
	}
	err = app.storage.Client.Set(app.ctx, shortURL, originURL, 0).Err()
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

func (app *application) FindOriginalURL(shortURL string) (string, error) {
	originURL, err := app.storage.Client.Get(app.ctx, shortURL).Result()
	if errors.Is(err, redis.Nil) {
		return "", err
	} else if err != nil {
		return "", err
	}
	return originURL, nil
}
