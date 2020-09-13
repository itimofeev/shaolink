package pg

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-pg/sharding/v8"
	"golang.org/x/exp/rand"

	"github.com/go-pg/pg/v10"

	"github.com/itimofeev/shaolink/internal/model"
)

// pow(len(alpha), 8) = 218,340,105,584,896 possibly saved links seems enough
const keyLength = 8
const shards = 10 // 0..9 logical shards, single symbol in key string

type Store struct {
	db      *pg.DB
	cluster *sharding.Cluster
}

func NewStore(connectString string) (*Store, error) {
	opts, err := pg.ParseURL(connectString)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	dbs := []*pg.DB{db} // list of physical PostgreSQL servers
	// Create cluster with 1 physical server and 2 logical shards.
	cluster := sharding.NewCluster(dbs, shards)

	// Create database schema for our logical shards.
	for i := 0; i < shards; i++ {
		if err := createShard(cluster.Shard(int64(i))); err != nil {
			panic(err)
		}
	}

	return &Store{
		db:      db,
		cluster: cluster,
	}, nil
}

func (s *Store) Save(original string) (*model.Shorten, error) {
	parsed, err := url.Parse(original)
	if err != nil {
		return nil, err
	}

	return s.SaveURL(parsed)
}

func (s *Store) SaveURL(originalURL *url.URL) (*model.Shorten, error) {
	shardNumber := rand.Intn(shards)

	toSave := &model.Shorten{
		Key:         fmt.Sprintf("%d%s", shardNumber, randomString(keyLength)),
		Original:    originalURL.String(),
		CreatedAt:   time.Now(),
		ShardNumber: shardNumber,
	}

	var err error
	for i := 0; i < 10; i++ {
		toSave.Key = fmt.Sprintf("%d%s", shardNumber, randomString(keyLength))
		err = trySave(s.cluster.Shard(int64(shardNumber)), toSave)
		if err == nil {
			break
		}

		pgErr, ok := err.(pg.Error)
		if !ok || !pgErr.IntegrityViolation() {
			return nil, err
		}
	}

	return toSave, err
}

func trySave(db *pg.DB, sh *model.Shorten) error {
	_, err := db.Model(sh).Insert()
	return err
}

func (s *Store) GetByKey(key string) (*model.Shorten, error) {
	shardNumberFromKey, err := strconv.Atoi(key[:1])
	if err != nil {
		return nil, err
	}

	loaded := &model.Shorten{
		Key: key,
	}
	return loaded, s.cluster.Shard(int64(shardNumberFromKey)).Model(loaded).WherePK().Select()
}

func (s *Store) LoadAll() ([]*model.Shorten, error) {
	mu := sync.Mutex{}
	all := make([]*model.Shorten, 0)

	return all, s.cluster.ForEachShard(func(shard *pg.DB) error {
		var ofShard []*model.Shorten
		if err := shard.Model(&ofShard).Select(); err != nil {
			return err
		}

		mu.Lock()
		all = append(all, ofShard...)
		mu.Unlock()
		return nil
	})
}

// createShard creates database schema for a given shard.
func createShard(shard *pg.DB) error {
	queries := []string{
		`CREATE SCHEMA IF NOT EXISTS ?SHARD`,
		`CREATE TABLE IF NOT EXISTS ?SHARD.shortens (key VARCHAR(16) NOT NULL PRIMARY KEY, original VARCHAR(2048) NOT NULL, created_at TIMESTAMPTZ NOT NULL, shard_number INTEGER NOT NULL)`, // nolint:lll // it's ok for SQL
	}

	for _, q := range queries {
		_, err := shard.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
