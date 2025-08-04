package mysql

import (
	"fmt"
	"sync"

	"github.com/pachirode/iam_study/internal/apiserver/store"
	"github.com/pachirode/iam_study/pkg/db"
	"github.com/pachirode/iam_study/pkg/errors"
	"gorm.io/gorm"

	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
)

type dataStore struct {
	db *gorm.DB
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

func (ds *dataStore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *dataStore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

func GetMySQLFactoryOr(opts *genericOptions.MySQLOptions) (store.Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("Failed to get mysql store factory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
		}
		dbIns, err = db.New(options)
		mysqlFactory = &dataStore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("Failed to get mysql store factory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}

func cleanDatabase(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&v1.User{}); err != nil {
		return errors.Wrap(err, "drop user table failed")
	}

	return nil
}

func migrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&v1.User{}); err != nil {
		return errors.Wrap(err, "Migrate user model failed")
	}

	return nil
}

func resetDatabase(db *gorm.DB) error {
	if err := cleanDatabase(db); err != nil {
		return err
	}

	if err := migrateDatabase(db); err != nil {
		return err
	}

	return nil
}
