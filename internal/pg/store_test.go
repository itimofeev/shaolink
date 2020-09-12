package pg

import (
	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StoreTestSuite struct {
	suite.Suite
	store *Store
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (s *StoreTestSuite) SetupSuite() {
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

func (s *StoreTestSuite) TestErrOnLoadNotExisted() {
	_, err := s.store.GetByKey("-1")
	s.Require().EqualError(err, pg.ErrNoRows.Error())
}
