package conn

import (
	"github.com/go-redis/redis"

	"github.com/jadoint/micro/db"
)

// Clients connections
type Clients struct {
	DB    *db.ClientDB
	Cache *redis.Client
}
