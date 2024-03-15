package domain

type ItemState int64

const (
	ItemStateUnspecified ItemState = iota
	ItemStateListed
	ItemStateSold
)

type Item struct {
	ID                   string    `json:"item_id"`
	Name                 string    `json:"name"`
	State                ItemState `json:"state"`
	Price                int64     `json:"price"`
	NFTID                string    `json:"nft_id"`
	NFTAddressID         string    `json:"nft_address_id"`
	SmartContractAddress string    `json:"smart_contract_address"`
}
