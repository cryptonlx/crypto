package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Client struct {
	serverUrl string
}

type WalletBalanceResponse struct {
	Error string `json:"error"`
}

func (c *Client) GetWalletBalance(_userId int64) (WalletBalanceResponse, int, error) {
	userId := strconv.Itoa(int(_userId))
	baseUrl := c.serverUrl + "/user/" + userId + "/balance"

	fullURL := baseUrl

	resp, err := http.Get(fullURL)
	if err != nil {
		return WalletBalanceResponse{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WalletBalanceResponse{}, 0, err
	}

	var b WalletBalanceResponse
	err = json.Unmarshal(body, &b)
	if err != nil {
		return WalletBalanceResponse{}, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return b, resp.StatusCode, fmt.Errorf("%s", b.Error)
	}
	return b, 0, nil
}

func NewClient(serverUrl string) (*Client, error) {
	return &Client{
		serverUrl: serverUrl,
	}, nil
}
