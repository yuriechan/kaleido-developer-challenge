package item

import (
	"backend/internal/domain"
	"context"
	"fmt"
)

type blockchainClient interface {
	CreatePool(ctx context.Context) error
	MintToken(ctx context.Context) error
	TransferToken(ctx context.Context) error
	DeploySmartContract(ctx context.Context, item *domain.Item) error
	EnableNFTForSale(ctx context.Context, contractAddress string) error
}

type dbClient interface {
	GetItemByID(ctx context.Context, id string) (*domain.Item, error)
	CreateItem(ctx context.Context, item *domain.Item) error
	UpdateItem(ctx context.Context, item *domain.Item) error
	CreateOrUpdateItem(ctx context.Context, item *domain.Item) error
}

type Service struct {
	blockchainClient blockchainClient
	dbClient         dbClient
}

func New(blockchainClient blockchainClient, dbClient dbClient) *Service {
	return &Service{
		blockchainClient: blockchainClient,
		dbClient:         dbClient,
	}
}

func (s *Service) GetItem(ctx context.Context, id string) (*domain.Item, error) {
	resp, err := s.dbClient.GetItemByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.dbClient.GetItemByID: %w", err)
	}
	return resp, nil
}

// ListItem can be used for listing a new item or/and re-listing an existing item
func (s *Service) ListItem(ctx context.Context, item *domain.Item) error {
	// TODO: Make calls to dbClient and blockchainClient as a transaction
	if err := s.dbClient.CreateOrUpdateItem(ctx, item); err != nil {
		return fmt.Errorf("s.dbClient.CreateOrUpdateItem: %w", err)
	}
	if err := s.blockchainClient.DeploySmartContract(ctx, item); err != nil {
		return fmt.Errorf("s.blockchainClient.DeploySmartContract: %w", err)
	}
	if err := s.blockchainClient.EnableNFTForSale(ctx, item.SmartContractAddress); err != nil {
		return fmt.Errorf("s.blockchainClient.EnableNFTForSale: %w", err)
	}
	return nil
}

func (s *Service) PurchaseItem(ctx context.Context, itemID string) error {
	resp, err := s.dbClient.GetItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("s.dbClient.GetItemByID: %w", err)
	}
	if resp.State == domain.ItemStateSold {
		return fmt.Errorf("item cannot be called when state is ItemStateSold")
	}

	resp.State = domain.ItemStateSold
	if err := s.dbClient.UpdateItem(ctx, resp); err != nil {
		return fmt.Errorf("s.dbClient.UpdateItem: %w", err)
	}
	if err := s.blockchainClient.TransferToken(ctx); err != nil {
		return fmt.Errorf("s.blockchainClient.TransferToken: %w", err)
	}
	return nil
}
