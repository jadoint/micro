package conn

import (
	"github.com/go-redis/redis/v7"

	"github.com/jadoint/micro/pkg/db"
)

// Clients connections
type Clients struct {
	DB    *db.ClientDB
	Cache *redis.Client
}
