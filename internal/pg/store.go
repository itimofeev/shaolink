package pg

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"strings"
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


func createSchema(db *pg.DB) error {
	for _, entity := range []interface{}{
		(*model.DeviceVMSConnect)(nil),
		(*model.DeviceTunnel)(nil),
		(*model.DeviceDirect)(nil),
		(*model.DeviceDirectPorts)(nil),
	} {
		err := db.CreateTable(entity, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	migs := getMigrations()
	_, _, err := migs.Run(db, "init")
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

	oldVersion, newVersion, err := migs.Run(db, "up")
	util.LogWithError(util.Log.WithField("oldVersion", oldVersion).WithField("newVersion", newVersion), err, "db schema migrated")
	if err != nil {
		return err
	}

	return nil
}
