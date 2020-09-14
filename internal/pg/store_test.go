package pg

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	store *Store
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (s *StoreTestSuite) SetupSuite() {
	rand.Seed(time.Now().Unix())
	store, err := NewStore("postgresql://postgres:@db:5432/postgres?sslmode=disable")
	s.Require().NoError(err)

	s.store = store
}

func (s *StoreTestSuite) TestConnectDB() {
	saved, err := s.store.Save("http://hello.world")
	s.Require().NoError(err)
	s.Require().Equal("http://hello.world", saved.Original)

	loaded, err := s.store.GetByKey(saved.Key)
	s.Require().NoError(err)

	s.Require().Equal(saved.Key, loaded.Key)
	s.Require().Equal(saved.Original, loaded.Original)
}

func (s *StoreTestSuite) TestSaveMultipleRecords() {
	for i := 0; i < 100; i++ {
		_, err := s.store.Save(fmt.Sprintf("http://%s.com/hi/there", randomString(10)))
		s.Require().NoError(err)
	}

	shortens, err := s.store.LoadAll()
	s.Require().NoError(err)

	distribution := make(map[int]int) // count of records in each shard
	for _, sh := range shortens {
		distribution[sh.ShardNumber]++
	}

	fmt.Println("distribution of records by shard number", distribution)
}

func (s *StoreTestSuite) TestErrOnLoadNotExisted() {
	_, err := s.store.GetByKey("1a")
	s.Require().EqualError(err, pg.ErrNoRows.Error())
}
