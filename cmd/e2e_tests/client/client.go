package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type Client struct {
	serverUrl  string
	httpClient *http.Client
}

func NewClient(serverUrl string) (*Client, error) {
	return &Client{
		serverUrl: serverUrl,
		httpClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       3 * time.Second,
		},
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

type WalletBalanceResponseBody = ResponseBody[WalletBalanceResponseData]

func (c *Client) Wallets(username string) (_responseBody WalletBalanceResponseBody, _httpStatusCode int, _clientError error) {
	if username == "" {
		return WalletBalanceResponseBody{}, 0, fmt.Errorf("malformedclient request. abort sending")
	}
	baseUrl := c.serverUrl + "/user/" + username + "/wallets"
	return httpGet[WalletBalanceResponseBody](c.httpClient, baseUrl, nil)
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

	return httpPost[CreateUserResponseBody](c.httpClient, baseUrl, requestBody, nil)
}

type TransactionMetaData struct {
	SourceWalletId *int64  `json:"source_wallet_id" example:"1"`
	Amount         *string `json:"amount" example:"1"`
}

type Transaction struct {
	Ledgers []Ledger `json:"ledgers"`

	Id          int64               `json:"id" example:"1"`
	RequestorId int64               `json:"requestor_id" example:"1"`
	Nonce       int64               `json:"nonce"`
	Status      string              `json:"status"`
	Operation   string              `json:"operation"`
	CreatedAt   time.Time           `json:"created_at"`
	MetaData    TransactionMetaData `json:"metadata"`
}

type TransactionResponseData struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionResponseBody = ResponseBody[TransactionResponseData]

func (c *Client) Transactions(username string) (TransactionResponseBody, int, error) {
	if username == "" {
		return TransactionResponseBody{}, 0, fmt.Errorf("malformedclient request. abort sending")
	}
	baseUrl := c.serverUrl + "/user/" + username + "/transactions"
	return httpGet[TransactionResponseBody](c.httpClient, baseUrl, nil)
}

type User struct {
	Username string `json:"username" example:"tester_123"`
	Id       int64  `json:"id" example:"1"`
}

type Ledger struct {
	Id        int64     `json:"id" example:"1"`
	WalletId  int64     `json:"wallet_id" example:"1"`
	EntryType string    `json:"entry_type" example:"credit,debit"`
	Amount    string    `json:"amount" example:"40.22"`
	CreatedAt time.Time `json:"created_at"`
	Balance   string    `json:"balance" example:"2.234"`
}

type DepositResponseData struct {
	Transaction `json:"transaction"`
}

type DepositResponseBody = ResponseBody[DepositResponseData]

func (c *Client) Deposit(username string, walletId int64, amount decimal.Decimal) (DepositResponseBody, int, error) {
	baseUrl := c.serverUrl + fmt.Sprintf("/wallet/%d/deposit", walletId)
	requestBody := map[string]interface{}{
		"amount": amount.String(),
		"nonce":  time.Now().UnixMilli(),
	}

	return httpPost[DepositResponseBody](c.httpClient, baseUrl, requestBody, []string{username, ""})
}

type WithdrawResponseData struct {
	Transaction `json:"transaction"`
}

type WithdrawResponseBody = ResponseBody[WithdrawResponseData]

func (c *Client) Withdraw(username string, walletId int64, amount decimal.Decimal) (WithdrawResponseBody, int, error) {
	baseUrl := c.serverUrl + fmt.Sprintf("/wallet/%d/withdraw", walletId)
	requestBody := map[string]interface{}{
		"amount": amount.String(),
		"nonce":  time.Now().UnixMilli(),
	}

	return httpPost[WithdrawResponseBody](c.httpClient, baseUrl, requestBody, []string{username, ""})
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
	return httpPost[CreateWalletResponseBody](c.httpClient, baseUrl, requestBody, nil)
}

type TransferResponseData struct {
	Transaction `json:"transaction"`
}

type TransferResponseBody = ResponseBody[TransferResponseData]

func (c *Client) Transfer(username string, fromWalletId int64, toWalletId int64, amount decimal.Decimal) (TransferResponseBody, int, error) {
	baseUrl := c.serverUrl + fmt.Sprintf("/wallet/%d/transfer", fromWalletId)
	requestBody := map[string]interface{}{
		"destination_wallet_id": toWalletId,
		"amount":                amount.String(),
		"nonce":                 time.Now().UnixMilli(),
	}
	return httpPost[TransferResponseBody](c.httpClient, baseUrl, requestBody, []string{username, ""})
}

type ResponseBody[T any] struct {
	Error *string `json:"error"`
	Data  T       `json:"data"`
}
