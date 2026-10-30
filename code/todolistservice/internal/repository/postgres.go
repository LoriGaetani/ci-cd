package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"example.com/todolistservice/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed schema.sql
var schema string

type Postgres struct {
	db *sql.DB
}

func Open(ctx context.Context, databaseURL string) (*Postgres, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := waitForDatabase(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize schema: %w", err)
	}
	return &Postgres{db: db}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) Ready(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *Postgres) List(ctx context.Context) ([]model.Item, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, description FROM todo_item ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Item, 0)
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (p *Postgres) Create(ctx context.Context, description string) (model.Item, error) {
	item := model.Item{Description: description}
	err := p.db.QueryRowContext(ctx,
		"INSERT INTO todo_item (description) VALUES ($1) RETURNING id", description).Scan(&item.ID)
	return item, err
}

func waitForDatabase(ctx context.Context, db *sql.DB) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	for {
		if err := db.PingContext(timeoutCtx); err == nil {
			return nil
		} else if timeoutCtx.Err() != nil {
			return fmt.Errorf("database connection timed out: %w", err)
		}
		select {
		case <-timeoutCtx.Done():
			return timeoutCtx.Err()
		case <-time.After(time.Second):
		}
	}
}
