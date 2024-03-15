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
)

type Client struct {
	httpClient *http.Client
	port       uidToHTTPPort
}

type uidToHTTPPort map[string]*url.URL

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
	approveTokenPath      = "tokens/approvals"
	applicationJsonHeader = "application/json"

	deployContractPath = "contracts/deploy"
	getTransactionPath = "transactions"

	buyNFTPath       = "apis/marketplace/invoke/buyNFT"
	nftPoolName      = "kaleido"
	nftPoolID        = "0xd9d2f32fecdbcaa40b48b03132dc1023fa63d171"
	nftDefaultAmount = "1"
	nftType          = "nonfungible"
	defaultUserID    = "1"
)

type createPoolRequest struct {
	PoolName string `json:"name"`
	PoolType string `json:"type"`
}

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
	Tx string `json:"tx"`
}

func (c *Client) DeploySmartContract(ctx context.Context, item *domain.Item) (string, error) {
	var input []string
	// TODO: Add logic when NFT data is empty
	input = append(input, nftPoolID, item.NFTID, strconv.Itoa(int(item.Price)))
	req := deploySmartContractRequest{
		Contract:   contracts.GetMarketplaceBin(),
		Definition: contracts.GetMarketplaceABI(),
		Input:      input,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("json.Marshal type deploySmartContractRequest: %w", err)
	}

	u := c.port[utils.FromContext(ctx)].JoinPath(deployContractPath)
	fmt.Println(u.String())
	resp, err := c.httpClient.Post(u.String(), applicationJsonHeader, bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("c.httpClient.Post to (%s): %w", u.String(), err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("io.ReadAll on response from (%s): %w", u.String(), err)
	}

	var res deploySmartContractResponse
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return "", fmt.Errorf("json.Unmarshal on response from (%s): %w", u.String(), err)
	}
	return res.Tx, nil
}

type getTransactionStatusResp struct {
	Details []struct {
		Info struct {
			ContractLocation struct {
				Address string `json:"address"`
			} `json:"contractLocation"`
		} `json:"info"`
	} `json:"details"`
}

func (c *Client) GetSmartContractLocation(ctx context.Context, trxID string) (string, error) {
	u := c.port[utils.FromContext(ctx)].JoinPath(getTransactionPath).JoinPath(trxID).JoinPath("status")
	fmt.Println(u.String())
	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("c.httpClient.Post to (%s): %w", u.String(), err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("io.ReadAll on response from (%s): %w", u.String(), err)
	}

	var res getTransactionStatusResp
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return "", fmt.Errorf("json.Unmarshal on response from (%s): %w", u.String(), err)
	}
	return res.Details[0].Info.ContractLocation.Address, nil
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

func (c *Client) ApproveTokenTransfer(ctx context.Context, item *domain.Item) error {
	req := approveTokenTransferRequest{
		Operator: item.SmartContractAddress,
		Config: struct {
			TokenID string `json:"tokenIndex"`
		}{
			TokenID: item.NFTID,
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

type mintTokenRequest struct {
	Pool   string `json:"pool"`
	Amount string `json:"amount"`
}

type mintTokenResponse struct {
	TokenIndex string `json:"tokenIndex"`
}

func (c *Client) MintToken(ctx context.Context) (string, error) {
	req := mintTokenRequest{
		Pool:   nftPoolName,
		Amount: nftDefaultAmount,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("json.Marshal type mintTokenRequest: %w", err)
	}
	p := c.port[utils.FromContext(ctx)].JoinPath(mintTokenPath)
	resp, err := c.httpClient.Post(p.String(), applicationJsonHeader, bytes.NewBuffer(b))
	if err != nil {
		return "", fmt.Errorf("c.httpClient.Post to (%s): %w", p.String(), err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	fmt.Println("--------------- after mint")
	fmt.Println(string(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("io.ReadAll on response from (%s): %w", p.String(), err)
	}
	var res mintTokenResponse
	if err = json.Unmarshal(bodyBytes, &res); err != nil {
		return "", fmt.Errorf("json.Unmarshal on response from (%s): %w", p.String(), err)
	}
	fmt.Println(res.TokenIndex)
	return res.TokenIndex, nil
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
