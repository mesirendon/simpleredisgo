package datasource

import (
	"os"

	"github.com/go-redis/redis/v8"
)

type Database int

const (
	DefaultDatabase Database = iota
	RulesDatabase
)

// Creates a new client to the specified database
func CreateClient(db Database) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("DB_ADDR"),
		DB:   int(db),
	})
}
