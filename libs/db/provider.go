package db

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

var (
	dbMtx   sync.RWMutex
	dbCache = make(map[string]*sqlx.DB)
)

/*
'postgres' driver is registered from lib/pq/conn.go init
*/
const driver string = "postgres"

// Provider ...
type Provider interface {
	Get(ctx context.Context, cfg Config) (*sqlx.DB, error)
}

// provider ...
type provider struct {
	dbMtx   sync.RWMutex
	dbCache map[string]*sqlx.DB
}

// NewProvider ...
func NewProvider() Provider {
	return &provider{}
}

// Config ...
type Config struct {
	ConnectionString             *string
	Driver                       *string
	MaxOpenConnections           *int
	MaxIdleConnections           *int
	ConnectionMaxLifetimeMinutes *int
}

// Get ...
func (p *provider) Get(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	sqltrace.Register(driver, pq.Driver{})

	if cfg.ConnectionString == nil || len(*cfg.ConnectionString) < 1 {
		return nil, errors.New("no connection string passed")
	}

	var db *sqlx.DB

	dbMtx.RLock()
	db = dbCache[*cfg.ConnectionString]
	dbMtx.RUnlock()

	if db != nil {
		return db, nil
	}

	dbMtx.Lock()
	defer dbMtx.Unlock()

	db = dbCache[*cfg.ConnectionString]
	if db != nil {
		return db, nil
	}

	dbs, err := sqltrace.Open(driver, *cfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	if cfg.MaxOpenConnections != nil {
		dbs.SetMaxOpenConns(*cfg.MaxOpenConnections)
	}
	if cfg.MaxIdleConnections != nil {
		dbs.SetMaxIdleConns(*cfg.MaxIdleConnections)
	}
	if cfg.ConnectionMaxLifetimeMinutes != nil {
		dbs.SetConnMaxLifetime(time.Duration(*cfg.MaxOpenConnections) * time.Minute)
	}

	db = sqlx.NewDb(dbs, driver)

	dbCache[*cfg.ConnectionString] = db
	return db, nil
}
