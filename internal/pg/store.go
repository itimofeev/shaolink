package pg

import (
	"net/url"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/rs/xid"

	"github.com/itimofeev/shaolink/internal/model"
)

type Store struct {
	db *pg.DB
}

// "postgresql://postgres:@db:5432/postgres?sslmode=disable"
func NewStore(connectString string) (*Store, error) {
	opts, err := pg.ParseURL(connectString)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	return &Store{
		db: db,
	}, createSchema(db)
}

func (s *Store) Save(original string) (*model.Shorten, error) {
	parsed, err := url.Parse(original)
	if err != nil {
		return nil, err
	}

	return s.SaveURL(parsed)
}

func (s *Store) SaveURL(originalURL *url.URL) (*model.Shorten, error) {
	toSave := &model.Shorten{
		Key:       xid.New().String(),
		Original:  originalURL.String(),
		CreatedAt: time.Now(),
	}

	_, err := s.db.Model(toSave).Insert()
	return toSave, err
}

func (s *Store) GetByKey(key string) (*model.Shorten, error) {
	loaded := &model.Shorten{
		Key: key,
	}
	return loaded, s.db.Model(loaded).WherePK().Select()
}

func createSchema(db *pg.DB) error {
	for _, entity := range []interface{}{
		(*model.Shorten)(nil),
	} {
		err := db.Model(entity).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
