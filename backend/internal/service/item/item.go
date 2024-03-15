package item

import (
	"backend/internal/domain"
	"context"
	"fmt"
	"time"
)

type fireflyClient interface {
	DeploySmartContract(ctx context.Context, item *domain.Item) (string, error)
	ApproveTokenTransfer(ctx context.Context, item *domain.Item) error
	BuyNFT(ctx context.Context, contractAddress string) error
	MintToken(ctx context.Context) (string, error)
	GetSmartContractLocation(ctx context.Context, trxID string) (string, error)
}

type dbClient interface {
	GetItemByID(ctx context.Context, id string) (*domain.Item, error)
	CreateItem(ctx context.Context, item *domain.Item) error
	UpdateItem(ctx context.Context, item *domain.Item) error
	CreateOrUpdateItem(ctx context.Context, item *domain.Item) error
}

type Service struct {
	fireflyClient fireflyClient
	dbClient      dbClient
}

func New(fireflyClient fireflyClient, dbClient dbClient) *Service {
	return &Service{
		fireflyClient: fireflyClient,
		dbClient:      dbClient,
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
	if err := s.dbClient.CreateOrUpdateItem(ctx, item); err != nil {
		return fmt.Errorf("ListItem: s.dbClient.CreateOrUpdateItem: %w", err)
	}

	trxID, err := s.fireflyClient.DeploySmartContract(ctx, item)
	if err != nil {
		return fmt.Errorf("ListItem: s.fireflyClient.DeploySmartContract: %w", err)
	}

	// TODO: Use event listener
	time.Sleep(time.Second * 5)
	clocation, err := s.fireflyClient.GetSmartContractLocation(ctx, trxID)
	if err != nil {
		return fmt.Errorf("ListItem: s.fireflyClient.GetSmartContractLocation: %w", err)
	}

	item.SmartContractAddress = clocation
	if err := s.dbClient.UpdateItem(ctx, item); err != nil {
		return fmt.Errorf("ListItem: s.dbClient.UpdateItem: %w", err)
	}

	if err := s.fireflyClient.ApproveTokenTransfer(ctx, item); err != nil {
		return fmt.Errorf("s.fireflyClient.ApproveTokenTransfer: %w", err)
	}
	return nil
}

func (s *Service) PurchaseItem(ctx context.Context, item *domain.Item) error {
	resp, err := s.dbClient.GetItemByID(ctx, item.ID)
	if err != nil {
		return fmt.Errorf("s.dbClient.GetItemByID: %w", err)
	}
	if resp.State == domain.ItemStateSold {
		return fmt.Errorf("item cannot be called when state is ItemStateSold")
	}

	if err := s.fireflyClient.BuyNFT(ctx, resp.SmartContractAddress); err != nil {
		return fmt.Errorf("s.fireflyClient.TransferToken: %w", err)
	}
	resp.State = domain.ItemStateSold
	if err := s.dbClient.UpdateItem(ctx, resp); err != nil {
		return fmt.Errorf("s.dbClient.UpdateItem: %w", err)
	}
	return nil
}
