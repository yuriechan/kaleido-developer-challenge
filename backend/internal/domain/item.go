package domain

type ItemState int64

const (
	ItemStateUnspecified ItemState = iota
	ItemStateListed
	ItemStateSold
)

type Item struct {
	ID                   string `json:"item_id"`
	Name                 string `json:"item_name"`
	State                ItemState
	Price                int64  `json:"item_price"`
	NFTID                string `json:"nft_id"`
	SmartContractAddress string `json:"smart_contract_address"`
}
