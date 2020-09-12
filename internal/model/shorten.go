package model

import "time"

type Shorten struct {
	Key       string    `pg:",pk,notnull"`
	Original  string    `pg:",notnull"`
	CreatedAt time.Time `pg:",notnull"`
}
