package blockchain

import (
	"backend/contracts"
	"backend/internal/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	url        *url.URL
	httpClient *http.Client
}

func New(u *url.URL, httpClient *http.Client) *Client {
	return &Client{
		url:        u,
		httpClient: httpClient,
	}
}

const (
	createPoolPath        = "tokens/pools"
	mintTokenPath         = "tokens/mint"
	transferTokenPath     = "tokens/transfers"
	applicationJsonHeader = "application/json"

	deployContractPath = "contracts/deploy"

	createNFTPath = "apis/marketplace/invoke/createNFT"
)

// TODO: Call this endpoint after HTTP server is up and running, before accepting requests
func (c *Client) CreatePool(ctx context.Context) error {
	p := c.url.JoinPath(createPoolPath)
	body := []byte(`{
		"name": "kaleido",
		"type": "nonfungible"
	  }`)
	_, err := c.httpClient.Post(p.String(), applicationJsonHeader, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", p, err)
	}
	return nil
}

type deploySmartContractRequest struct {
	Contract   string                     `json:"contract"`
	Definition contracts.SmartContractABI `json:"definition"`
	Input      []string                   `json:"input"`
}

type deploySmartContractResponse struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Tx        string `json:"tx"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Plugin    string `json:"plugin"`
	Input     struct {
		Contract   string `json:"contract"`
		Definition []struct {
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
		} `json:"definition"`
		Input   []string `json:"input"`
		Key     string   `json:"key"`
		Options any      `json:"options"`
	} `json:"input"`
	Output struct {
		Headers struct {
			RequestID string `json:"requestId"`
			Type      string `json:"type"`
		} `json:"headers"`
		ContractLocation struct {
			Address string `json:"address"`
		} `json:"contractLocation"`
		ProtocolID      string `json:"protocolId"`
		TransactionHash string `json:"transactionHash"`
	} `json:"output"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (c *Client) DeploySmartContract(ctx context.Context, item *domain.Item) error {
	var input []string
	input = append(input, item.NFTAddressID, item.NFTID, strconv.Itoa(int(item.Price)))
	req := deploySmartContractRequest{
		Contract:   contracts.GetMarketplaceBin(),
		Definition: contracts.GetMarketplaceABI(),
		Input:      input,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("json.Marshal type deploySmartContractRequest: %w", err)
	}

	u := c.url.JoinPath(deployContractPath)
	fmt.Println(u.String())
	resp, err := c.httpClient.Post(u.String(), applicationJsonHeader, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", u.String(), err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
	if err != nil {
		return fmt.Errorf("io.ReadAll on response from (%s): %w", u.String(), err)
	}

	// var res deploySmartContractResponse
	// if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
	// 	return fmt.Errorf("json.NewDecoder on response from (%s): %w", u.String(), err)
	// }
	return nil
}

type enableNFTForSaleRequest struct {
	Location *location         `json:"location"`
	Key      string            `json:"key"`
	Input    map[string]string `json:"input"`
}

type location struct {
	ContractAddress string `json:"address"`
}

func (c *Client) EnableNFTForSale(ctx context.Context, contractAddress string) error {
	req := enableNFTForSaleRequest{
		Location: &location{
			ContractAddress: contractAddress,
		},
	}
	b, err := json.Marshal(req)
	fmt.Println(string(b))
	if err != nil {
		return fmt.Errorf("json.Marshal type enableNFTForSaleRequest: %w", err)
	}

	u := c.url.JoinPath(createNFTPath)
	fmt.Println(u.String())
	resp, err := c.httpClient.Post(u.String(), applicationJsonHeader, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", u.String(), err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
	if err != nil {
		return fmt.Errorf("io.ReadAll on response from (%s): %w", u.String(), err)
	}

	return nil
}

func (c *Client) MintToken(ctx context.Context) error {
	p := c.url.JoinPath(mintTokenPath)
	body := []byte(`{
		"pool": "kaleido",
		"amount": "1"
	  }`)
	_, err := c.httpClient.Post(p.String(), applicationJsonHeader, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", p, err)
	}
	// TODO: store IDs into DB
	return nil
}

func (c *Client) TransferToken(ctx context.Context) error {
	p := c.url.JoinPath(transferTokenPath)
	// TODO: recipient address should be dynamic
	body := []byte(`{
		pool: 'kaleido',
		to: '0x2c499018bb8bc8a56065c07daeff4bceec254928',
		tokenIndex: '1',
		amount: '1',
	  }`)
	_, err := c.httpClient.Post(p.String(), applicationJsonHeader, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", p, err)
	}
	// TODO: store IDs into DB
	return nil
}
