package repository

import (
	"github.com/go-redis/redis/v8"
)

type ProductsReddisRepository struct {
	db *redis.Client
}

