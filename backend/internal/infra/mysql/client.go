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
	fmt.Println("inside infra before query")
	fmt.Println(id)
	query := "SELECT id, item_name, item_state, item_price, nft_id, smart_contract_address FROM listing WHERE id = ?"
	if err := c.db.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.State, &item.Price, &item.NFTID, &item.SmartContractAddress); err != nil {
		if errors.As(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("c.db.QueryRow on (%s) with id (%s): %w", query, id, err)
	}
	fmt.Println("inside infra after query")
	fmt.Println(item.ID)
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
	updateQuery := "UPDATE listing SET item_name = ?, item_state = ?, item_price = ?, nft_id = ?, smart_contract_address = ? WHERE id = ?"
	if _, err := c.db.ExecContext(ctx, updateQuery, item.Name, item.State, item.Price, item.NFTID, item.SmartContractAddress, item.ID); err != nil {
		return fmt.Errorf("c.db.QueryRow on (%s) with id (%s): %w", updateQuery, item.ID, err)
	}
	return nil
}

func (c *Client) CreateOrUpdateItem(ctx context.Context, item *domain.Item) error {
	if item == nil {
		return fmt.Errorf("CreateOrUpdateItem called with nil item data")
	}

	selectQuery := "SELECT id FROM listing WHERE id = ?"
	var isCreated bool
	var id string
	if err := c.db.QueryRow(selectQuery, item.ID).Scan(&id); errors.Is(err, sql.ErrNoRows) {
		insertQuery := "INSERT INTO listing (id, item_name, item_state, item_price, smart_contract_address, nft_id) VALUES (?, ?, ?, ?, ?, ?)"
		item.ID = uuid.NewString()
		if _, err := c.db.ExecContext(ctx, insertQuery, item.ID, item.Name, item.State, item.Price, item.SmartContractAddress, item.NFTID); err != nil {
			return fmt.Errorf("c.db.ExecContext on (%s) with id (%s): %w", insertQuery, item.ID, err)
		}
		isCreated = true
	} else if err != nil {
		return fmt.Errorf("c.db.QueryRow on (%s) with id (%s): %w", selectQuery, item.ID, err)
	}

	if isCreated {
		return nil
	}

	updateQuery := "UPDATE listing SET item_name = ?, item_state = ?, item_price = ?, smart_contract_address = ?, nft_id = ? WHERE id = ?"
	if _, err := c.db.ExecContext(ctx, updateQuery, item.Name, item.State, item.Price, item.SmartContractAddress, item.NFTID, item.ID); err != nil {
		return fmt.Errorf("c.db.QueryRow on (%s) with id (%s): %w", updateQuery, item.ID, err)
	}
	return nil
}
