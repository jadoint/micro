package conn

import (
	"github.com/redis/go-redis/v9"

	"github.com/jadoint/micro/pkg/db"
)

// Clients connections
type Clients struct {
	DB    *db.ClientDB
	Cache *redis.Client
}
