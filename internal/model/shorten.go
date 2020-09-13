package model

import "time"

type Shorten struct {
	tableName string `pg:"?SHARD.shortens"` // nolint:structcheck,unused // it's for go-pg library

	Key         string    `pg:",pk,notnull"`
	Original    string    `pg:",notnull"`
	CreatedAt   time.Time `pg:",notnull"`
	ShardNumber int       `pg:",notnull,use_zero"`
}
