package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type FileEntity struct {
	Id        int64     `json:"id"`
	Filename  string    `json:"filename"`
	Filetype  string    `json:"filetype"`
	Path      string    `json:"path"`
	SizeBytes int       `json:"size_bytes"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	CreatedAt time.Time `json:"created_at"`
}

type FileEntityStore struct {
	db *DatabaseConnection
}

func NewFileEntityStore(database *DatabaseConnection) *FileEntityStore {
	return &FileEntityStore{
		db: database,
	}
}

func (c *FileEntityStore) Exists(id int64) (bool, error) {
	row := c.db.Connection.QueryRow(c.db.Context, `select exists(select 1 from "files" where id = $1)`, id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *FileEntityStore) Create(ctx context.Context, file FileEntity) (result *FileEntity, err error) {
	tx, err := c.db.Connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	result = &FileEntity{
		Filename:  file.Filename,
		Filetype:  file.Filetype,
		Path:      file.Path,
		SizeBytes: file.SizeBytes,
		Width:     file.Width,
		Height:    file.Height,
	}

	err = tx.QueryRow(c.db.Context, `
		insert into "files" (filename, filetype, path, width, height, size_bytes) values ($1, $2, $3, $4, $5, $6) 
		returning id, created_at`, file.Filename, file.Filetype, file.Path, file.Width, file.Height, file.SizeBytes).Scan(&result.Id, &result.CreatedAt)

	if err != nil {
		return result, fmt.Errorf("failed to insert file: %w", err)
	}

	return result, tx.Commit(ctx)
}
