package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	serverUrl string
}

func NewClient(serverUrl string) (*Client, error) {
	return &Client{
		serverUrl: serverUrl,
	}, nil
}

type Wallet struct {
	Id            int64  `json:"id"`
	UserAccountId int64  `json:"user_account_id"`
	Currency      string `json:"currency"`
	Balance       string `json:"balance"`
}

type WalletBalanceResponseData struct {
	Wallets []Wallet `json:"wallets"`
	User    `json:"user"`
}

type WalletBalanceResponseBody struct {
	Error *string                   `json:"error"`
	Data  WalletBalanceResponseData `json:"data"`
}

func (c *Client) GetWalletBalances(username string) (_responseBody WalletBalanceResponseBody, _httpStatusCode int, _clientError error) {
	if username == "" {
		return WalletBalanceResponseBody{}, 0, fmt.Errorf("malformedclient request. abort sending")
	}
	baseUrl := c.serverUrl + "/user/" + username + "/wallets"

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
	*CreatedUser `json:"user"`
}

type CreateUserResponseBody = ResponseBody[CreateUserResponseData]

func (c *Client) CreateUser(username string) (CreateUserResponseBody, int, error) {
	baseUrl := c.serverUrl + "/user"

	requestBody := map[string]interface{}{
		"username": username,
	}

	return httpPost[CreateUserResponseBody](baseUrl, requestBody)
}

type Transaction struct{}

type TransactionResponseBodyData struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionResponseBody struct {
	Error *string                     `json:"error"`
	Data  TransactionResponseBodyData `json:"data"`
}

func (c *Client) GetTransactionHistory(username string) (TransactionResponseBody, int, error) {
	if username == "" {
		return TransactionResponseBody{}, 0, fmt.Errorf("malformedclient request. abort sending")
	}
	baseUrl := c.serverUrl + "/user/" + username + "/transactions"

	resp, err := http.Get(baseUrl)
	if err != nil {
		return TransactionResponseBody{}, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	//log.Println(string(body))
	//log.Println(resp.StatusCode)
	if err != nil {
		return TransactionResponseBody{}, 0, err
	}

	var b TransactionResponseBody
	err = json.Unmarshal(body, &b)
	if err != nil {
		return TransactionResponseBody{}, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s", *b.Error)
	}
	return b, resp.StatusCode, err
}

type User struct {
	Username string `json:"username" example:"tester_123"`
	Id       int64  `json:"id" example:"1"`
}

type DepositResponseData struct {
	User `json:"user"`
}

type DepositResponseBody = ResponseBody[DepositResponseData]

func (c *Client) Deposit(walletId int64, amount int64) (DepositResponseBody, int, error) {
	baseUrl := c.serverUrl + fmt.Sprintf("/wallet/%d/deposit", walletId)
	requestBody := map[string]interface{}{
		"amount": strconv.Itoa(int(amount)),
		"nonce":  time.Now().Unix(),
	}

	return httpPost[DepositResponseBody](baseUrl, requestBody)
}

type ResponseBody[T any] struct {
	Error *string `json:"error"`
	Data  T       `json:"data"`
}

type CreateWalletResponseData struct {
	Wallet *struct {
		Id            int64  `json:"id"`
		UserAccountId int64  `json:"user_account_id"`
		Currency      string `json:"currency"`
		Balance       string `json:"balance"`
	} `json:"wallet"`
}

type CreateWalletResponseBody = ResponseBody[CreateWalletResponseData]

func (c *Client) CreateWallet(username string, currency string) (CreateWalletResponseBody, int, error) {
	baseUrl := c.serverUrl + "/wallet"
	requestBody := map[string]interface{}{
		"username": username,
		"currency": currency,
	}
	return httpPost[CreateWalletResponseBody](baseUrl, requestBody)
}
