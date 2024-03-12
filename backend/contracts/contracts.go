package contracts

import (
	_ "embed"
	"encoding/json"
)

//go:embed Marketplace.abi
var mpABIRaw []byte

var mpABI SmartContractABI

//go:embed Marketplace.bin
var mpBin string

type SmartContractABI []struct {
	Inputs []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs"`
	StateMutability string `json:"stateMutability,omitempty"`
	Type            string `json:"type"`
	Anonymous       bool   `json:"anonymous,omitempty"`
	Name            string `json:"name,omitempty"`
	Outputs         []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"outputs,omitempty"`
}

func init() {
	abis := make(SmartContractABI, 0, 1)
	if err := json.Unmarshal(mpABIRaw, &abis); err != nil {
		panic(err)
	}
	mpABI = abis
}

func GetMarketplaceABI() SmartContractABI {
	return mpABI
}

func GetMarketplaceBin() string {
	return mpBin
}
