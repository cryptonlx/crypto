package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	serverUrl string
}

type WalletBalanceResponseBody struct {
	Error *string `json:"error"`
}

func (c *Client) GetWalletBalance(username string) (WalletBalanceResponseBody, int, error) {
	if username == "" {
		return WalletBalanceResponseBody{}, 0, fmt.Errorf("malformedclient request. abort sending")
	}
	baseUrl := c.serverUrl + "/user/" + username + "/balance"

	resp, err := http.Get(baseUrl)
	if err != nil {
		return WalletBalanceResponseBody{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	//log.Println(string(body))
	//log.Println(resp.StatusCode)
	if err != nil {
		return WalletBalanceResponseBody{}, 0, err
	}

	var b WalletBalanceResponseBody
	err = json.Unmarshal(body, &b)
	if err != nil {
		return WalletBalanceResponseBody{}, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s", *b.Error)
	}
	return b, resp.StatusCode, err
}

type CreatedUser struct {
	Username string `json:"username" example:"tester_123"`
	Id       int64  `json:"id" example:"1"`
}

type CreateUserResponseData struct {
	CreatedUser `json:"user"`
}

type CreateUserResponseBody struct {
	Error *string                `json:"error"`
	Data  CreateUserResponseData `json:"data"`
}

func (c *Client) CreateUser(username string) (CreateUserResponseBody, int, error) {
	baseUrl := c.serverUrl + "/user"

	requestBody := map[string]interface{}{
		"username": username,
	}

	bb, err := json.Marshal(requestBody)
	if err != nil {
		return CreateUserResponseBody{}, 0, err
	}

	fullURL := baseUrl
	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(bb))
	if err != nil {
		return CreateUserResponseBody{}, 0, err
	}
	defer resp.Body.Close()

	//log.Printf("resp.StatusCode: %v\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CreateUserResponseBody{}, 0, err
	}

	//log.Printf("Body: %s\n", string(body))

	var b CreateUserResponseBody
	err = json.Unmarshal(body, &b)
	if err != nil {
		return CreateUserResponseBody{}, 0, fmt.Errorf("response body: %s %v", string(body), err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s", *b.Error)
	}
	return b, resp.StatusCode, err
}

func NewClient(serverUrl string) (*Client, error) {
	return &Client{
		serverUrl: serverUrl,
	}, nil
}
