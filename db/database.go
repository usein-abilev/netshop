package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// Represents an active connection to postgresql database
type DatabaseConnection struct {
	Connection       *pgx.Conn
	Context          context.Context
	ConnectionString string
}

func NewDatabaseConnection(ctx context.Context, connectionStr string) (*DatabaseConnection, error) {
	connection, err := pgx.Connect(ctx, connectionStr)
	if err != nil {
		return &DatabaseConnection{}, err
	}
	return &DatabaseConnection{
		Context:          ctx,
		Connection:       connection,
		ConnectionString: connectionStr,
	}, nil
}

func (db *DatabaseConnection) Close() {
	db.Connection.Close(db.Context)
}
