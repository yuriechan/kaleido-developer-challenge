package firefly

import (
	"backend/contracts"
	"backend/internal/domain"
	"backend/internal/utils"
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
	httpClient *http.Client
	port       uidToHTTPPort
}

func New(u1, u2, u3 *url.URL, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		// For now, we use fixed userID
		port: uidToHTTPPort{
			"1": u1,
			"2": u2,
			"3": u3,
		},
	}
}

const (
	createPoolPath        = "tokens/pools"
	mintTokenPath         = "tokens/mint"
	transferTokenPath     = "tokens/transfers"
	approveTokenPath      = "tokens/approvals"
	applicationJsonHeader = "application/json"

	deployContractPath = "contracts/deploy"

	createNFTPath = "apis/marketplace/invoke/createNFT"
	buyNFTPath    = "apis/marketplace/invoke/buyNFT"
	nftPoolName   = "kaleido"
	nftType       = "nonfungible"
	defaultUserID = "1"
)

type createPoolRequest struct {
	PoolName string `json:"name"`
	PoolType string `json:"type"`
}

type uidToHTTPPort map[string]*url.URL

func (c *Client) CreatePool(_ context.Context) error {
	// For now, we default to using the user ID 1's port number for token pool creation
	p := c.port[defaultUserID].JoinPath(createPoolPath)
	req := createPoolRequest{
		PoolName: nftPoolName,
		PoolType: nftType,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("json.Marshal type createPoolRequest: %w", err)
	}
	if _, err = c.httpClient.Post(p.String(), applicationJsonHeader, bytes.NewBuffer(b)); err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", p, err)
	}
	return nil
}

type deploySmartContractRequest struct {
	Contract       string                     `json:"contract"`
	Definition     contracts.SmartContractABI `json:"definition"`
	Input          []string                   `json:"input"`
	IdempotencyKey string                     `json:"idempotencyKey"`
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
	// TODO: Add logic when NFT data is empty
	// TODO: Need to confirm what is NFT Address ID on Firefly Web UI
	// TODO: Listen to events after call to /deploy to retrieve location address during run time
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

	u := c.port[utils.FromContext(ctx)].JoinPath(deployContractPath)
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

	//var res deploySmartContractResponse
	//if err := json.Unmarshal(bodyBytes, &res); err != nil {
	//	return fmt.Errorf("json.Unmarshal on response from (%s): %w", u.String(), err)
	//}
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

type approveTokenTransferRequest struct {
	Operator string `json:"operator"`
	Config   struct {
		TokenID string `json:"tokenIndex"`
	} `json:"config"`
	Pool string `json:"pool"`
}

func (c *Client) ApproveTokenTransfer(ctx context.Context, nft *domain.NFT) error {
	req := approveTokenTransferRequest{
		Operator: nft.SmartContractID,
		Config: struct {
			TokenID string `json:"tokenIndex"`
		}{
			TokenID: strconv.Itoa(int(nft.ID)),
		},
		Pool: nftPoolName,
	}
	b, err := json.Marshal(req)
	fmt.Println(string(b))
	if err != nil {
		return fmt.Errorf("json.Marshal type approveTokenTransferRequest: %w", err)
	}
	u := c.port[utils.FromContext(ctx)].JoinPath(approveTokenPath)
	fmt.Println(u.String())
	resp, err := c.httpClient.Post(u.String(), applicationJsonHeader, bytes.NewReader(b))
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
	if err != nil {
		return fmt.Errorf("io.ReadAll on response from (%s): %w", u.String(), err)
	}
	return nil
}

func (c *Client) MintToken(ctx context.Context, uid string) error {
	p := c.port[utils.FromContext(ctx)].JoinPath(mintTokenPath)
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

type buyNFTLocation struct {
	Location struct {
		ContractAddress string `json:"address"`
	} `json:"location"`
}

func (c *Client) BuyNFT(ctx context.Context, contractAddress string) error {
	req := buyNFTLocation{Location: location{ContractAddress: contractAddress}}
	b, err := json.Marshal(req)
	fmt.Println(string(b))
	if err != nil {
		return fmt.Errorf("json.Marshal type location: %w", err)
	}

	u := c.port[utils.FromContext(ctx)].JoinPath(buyNFTPath)
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
