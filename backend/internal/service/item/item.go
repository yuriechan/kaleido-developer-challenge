package item

import (
	"backend/internal/domain"
	"context"
	"fmt"
)

type fireflyClient interface {
	DeploySmartContract(ctx context.Context, item *domain.Item) error
	ApproveTokenTransfer(ctx context.Context, nft *domain.NFT) error
	BuyNFT(ctx context.Context, contractAddress string) error
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

func (s *Service) ApproveTransferTokenOnBehalfOfBuyer(ctx context.Context, nft *domain.NFT) error {
	if err := s.fireflyClient.ApproveTokenTransfer(ctx, nft); err != nil {
		return fmt.Errorf("s.fireflyClient.ApproveTokenTransfer: %w", err)
	}
	return nil
}

// ListItem can be used for listing a new item or/and re-listing an existing item
func (s *Service) ListItem(ctx context.Context, item *domain.Item) error {
	// TODO: Make calls to dbClient and fireflyClient as a transaction
	if err := s.dbClient.CreateOrUpdateItem(ctx, item); err != nil {
		return fmt.Errorf("ListItem: s.dbClient.CreateOrUpdateItem: %w", err)
	}
	if err := s.fireflyClient.DeploySmartContract(ctx, item); err != nil {
		return fmt.Errorf("ListItem: s.fireflyClient.DeploySmartContract: %w", err)
	}
	return nil
}

func (s *Service) PurchaseItem(ctx context.Context, item *domain.Item) error {
	//resp, err := s.dbClient.GetItemByID(ctx, item.ID)
	//if err != nil {
	//	return fmt.Errorf("s.dbClient.GetItemByID: %w", err)
	//}
	//if resp.State == domain.ItemStateSold {
	//	return fmt.Errorf("item cannot be called when state is ItemStateSold")
	//}
	//
	//resp.State = domain.ItemStateSold
	//if err := s.dbClient.UpdateItem(ctx, resp); err != nil {
	//	return fmt.Errorf("s.dbClient.UpdateItem: %w", err)
	//}
	if err := s.fireflyClient.BuyNFT(ctx, item.SmartContractAddress); err != nil {
		return fmt.Errorf("s.fireflyClient.TransferToken: %w", err)
	}
	return nil
}
