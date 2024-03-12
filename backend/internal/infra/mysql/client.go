package mysql

import (
	"backend/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Client struct {
	db *sql.DB
}

func New(db *sql.DB) *Client {
	return &Client{
		db: db,
	}
}

func (c *Client) GetItemByID(ctx context.Context, id string) (*domain.Item, error) {
	var item domain.Item
	var row *sql.Row
	query := "SELECT id, name FROM items WHERE id = ?"
	if row = c.db.QueryRow(query, id); row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("c.db.QueryRow on (%s) with id (%s): %w", query, id, row.Err())
	}
	if err := row.Scan(&item.ID, &item.Name); err != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("row.Scan on (%s) with id (%s): %w", query, id, err)
	}
	return &item, nil
}

func (c *Client) CreateItem(ctx context.Context, item *domain.Item) error {
	if item == nil {
		return fmt.Errorf("CreateItem called with nil item data")
	}

	item.ID = uuid.NewString()
	insertQuery := "INSERT INTO items (id, name, state) VALUES (?, ?, ?)"
	if _, err := c.db.ExecContext(ctx, insertQuery, item.ID, item.Name, item.State); err != nil {
		return fmt.Errorf("c.db.ExecContext on (%s) with id (%s): %w", insertQuery, item.ID, err)
	}
	return nil
}

func (c *Client) UpdateItem(ctx context.Context, item *domain.Item) error {
	if item == nil {
		return fmt.Errorf("UpdateItem called with nil item data")
	}
	updateQuery := "UPDATE items SET name = ?, state = ? WHERE id = ?"
	if _, err := c.db.ExecContext(ctx, updateQuery, item.Name, item.State, item.ID); err != nil {
		return fmt.Errorf("c.db.QueryRow on (%s) with id (%s): %w", updateQuery, item.ID, err)
	}
	return nil
}

func (c *Client) CreateOrUpdateItem(ctx context.Context, item *domain.Item) error {
	if item == nil {
		return fmt.Errorf("CreateOrUpdateItem called with nil item data")
	}

	selectQuery := "SELECT id FROM items WHERE id = ?"
	var isCreated bool
	var id string
	if err := c.db.QueryRow(selectQuery, item.ID).Scan(&id); errors.Is(err, sql.ErrNoRows) {
		insertQuery := "INSERT INTO items (id, name, state) VALUES (?, ?, ?)"
		item.ID = uuid.NewString()
		if _, err := c.db.ExecContext(ctx, insertQuery, item.ID, item.Name, item.State); err != nil {
			return fmt.Errorf("c.db.ExecContext on (%s) with id (%s): %w", insertQuery, item.ID, err)
		}
		isCreated = true
	} else if err != nil {
		return fmt.Errorf("c.db.QueryRow on (%s) with id (%s): %w", selectQuery, item.ID, err)
	}

	if isCreated {
		return nil
	}

	updateQuery := "UPDATE items SET name = ?, state = ? WHERE id = ?"
	if _, err := c.db.ExecContext(ctx, updateQuery, item.Name, item.State, item.ID); err != nil {
		return fmt.Errorf("c.db.QueryRow on (%s) with id (%s): %w", updateQuery, item.ID, err)
	}
	return nil
}
