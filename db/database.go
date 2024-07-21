package db

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Represents an active connection to postgresql database
type DatabaseConnection struct {
	Connection       *pgxpool.Pool
	Context          context.Context
	ConnectionString string
}

type DatabaseConnectionOptions struct {
	ConnectionURL     string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
	ConnectTimeout    time.Duration
}

func parseConfig(opts *DatabaseConnectionOptions) (*pgxpool.Config, error) {
	dbConfig, err := pgxpool.ParseConfig(opts.ConnectionURL)
	if err != nil {
		return nil, err
	}

	if opts.MaxConns != 0 {
		dbConfig.MaxConns = opts.MaxConns
	}
	if opts.MinConns != 0 {
		dbConfig.MinConns = opts.MinConns
	}
	if opts.MaxConnLifetime != 0 {
		dbConfig.MaxConnLifetime = opts.MaxConnLifetime
	}
	if opts.MaxConnIdleTime != 0 {
		dbConfig.MaxConnIdleTime = opts.MaxConnIdleTime
	}
	if opts.HealthCheckPeriod != 0 {
		dbConfig.HealthCheckPeriod = opts.HealthCheckPeriod
	}
	if opts.ConnectTimeout != 0 {
		dbConfig.ConnConfig.ConnectTimeout = opts.ConnectTimeout
	}

	return dbConfig, nil
}

func NewDatabaseConnection(ctx context.Context, opts *DatabaseConnectionOptions) (*DatabaseConnection, error) {
	dbConfig, err := parseConfig(opts)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	return &DatabaseConnection{
		Context:          ctx,
		Connection:       pool,
		ConnectionString: opts.ConnectionURL,
	}, nil
}

func (db *DatabaseConnection) Close() {
	db.Connection.Close()
}
