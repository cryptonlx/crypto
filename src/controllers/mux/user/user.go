package user

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cryptonlx/crypto/src/controllers/response_types"
	userservice "github.com/cryptonlx/crypto/src/service/user"

	"github.com/shopspring/decimal"
)

type Handlers struct {
	service *userservice.Service
}

func NewHandlers(service *userservice.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

type User struct {
	Id       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"user1"`
}

type Wallet struct {
	Id            int64  `json:"id" example:"1"`
	UserAccountId int64  `json:"user_account_id" example:"1"`
	Currency      string `json:"currency" example:"USD"`
	Balance       string `json:"balance" example:"10.000123"`
}

type GetWalletsResponseData struct {
	User    User     `json:"user"`
	Wallets []Wallet `json:"wallets"`
}

type GetWalletsResponseBody = Response[GetWalletsResponseData]

// Wallets godoc
// @Summary      Get balances of user's wallets.
// @Description  Get balances of user's wallets.
// @Tags         wallet
// @Accept       application/json
// @Produce      application/json
// @Param        user_id   					path      string  true  "username"
// @Success      200  {object}  GetWalletsResponseBody
// @Failure      500  {object}  ErrorResponseBody
// @Router       /user/{username}/wallets [get]
func (h Handlers) Wallets(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//basicAuthB64, _ := ctx.Value("BASIC_AUTH").(string)
	//principal, err := mustExtractUsernameFromBasicAuthValue(basicAuthB64)
	//if err != nil {
	//	response_types.ErrorNoBody(w, http.StatusForbidden, err)
	//	return
	//}
	//if principal != userName {
	//	response_types.ErrorNoBody(w, http.StatusForbidden, fmt.Errorf("principal not allowed to access resource"))
	//	return
	//}
	userName := r.PathValue("username")

	walletBalances, err := h.service.GetUserWalletBalanceByUserName(r.Context(), userName)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	wallets := make([]Wallet, 0, len(walletBalances.Wallets))
	for _, wallet := range walletBalances.Wallets {
		wallets = append(wallets, Wallet{
			Id:            wallet.Id,
			UserAccountId: wallet.UserAccountId,
			Currency:      wallet.Currency,
			Balance:       wallet.Balance.String(),
		})
	}

	response_types.OkJsonBody(w, GetWalletsResponseData{
		User:    User(walletBalances.User),
		Wallets: wallets,
	})
}

type CreateUserRequestBody struct {
	UserName string `json:"username"`
}

type CreatedUser struct {
	Id       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"user1"`
}
type CreateUserResponseData struct {
	User *CreatedUser `json:"user"`
}

type CreateUserResponse = Response[CreateUserResponseData]

// CreateUser Create godoc
// @Summary      Create a new user.
// @Description  Create a new user.
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        request body CreateUserRequestBody true "Create User Request Body"
// @Success      200  {object}  CreateUserResponse
// @Failure      500  {object}  ErrorResponseBody
// @Router       /user [post]
func (h Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	form := &CreateUserRequestBody{}
	json.NewDecoder(r.Body).Decode(form)
	if form.UserName == "" {
		response_types.ErrorNoBody(w, http.StatusBadRequest, errors.New("user name is required"))
	}

	user, err := h.service.CreateUser(r.Context(), form.UserName)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	c := CreatedUser(user)
	response_types.OkJsonBody(w, CreateUserResponseData{User: &c})
}

type UserTransactions struct {
	User         User
	Transactions []Transaction
}

type TransactionsResponseData struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionsResponseBody = Response[TransactionsResponseData]

// Transactions godoc
// @Summary      Get transactions of user's wallets.
// @Description  Get transactions of user's wallets.
// @Tags         wallet
// @Accept       application/json
// @Produce      application/json
// @Param        user_id   					path      string  true  "username"
// @Success      200  {object}  TransactionsResponseBody
// @Failure      500  {object}  ErrorResponseBody
// @Router       /user/{username}/transactions [get]
func (h Handlers) Transactions(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("username")

	transactionLedgers, err := h.service.GetUserTransactionsByUserName(r.Context(), userName)
	if err != nil {

		log.Printf("%v\n", err)
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	transactions := make([]Transaction, 0, len(transactionLedgers))
	for _, transaction := range transactionLedgers {
		ledgers := make([]Ledger, 0, len(transaction.Ledgers))
		for _, ledger := range transaction.Ledgers {
			ledgers = append(ledgers, Ledger{
				Id:            ledger.Id,
				WalletId:      ledger.WalletId,
				TransactionId: ledger.TransactionId,
				EntryType:     ledger.EntryType,
				Amount:        ledger.Amount.String(),
				CreatedAt:     ledger.CreatedAt,
				Balance:       ledger.Balance.String(),
			})
		}

		transactions = append(transactions, Transaction{
			Ledgers: ledgers,
			Id:      transaction.Transaction.Id,
			//UserAccountId: transaction.UserAccountId,
			//Nonce:         transaction.Nonce,
			//Status:        transaction.Status,
			//Operation:     transaction.Operation,
			//CreatedAt:     transaction.CreatedAt,
		})
	}
	response_types.OkJsonBody(w, TransactionsResponseData{
		Transactions: transactions,
	})
}

type CreateWalletRequestBody struct {
	UserName string `json:"username"`
	Currency string `json:"currency"`
}

type CreatedWallet struct {
	Id            int64  `json:"id" example:"1"`
	UserAccountId int64  `json:"user_account_id" example:"1"`
	Balance       string `json:"balance" example:"user1"`
	Currency      string `json:"currency" example:"USD"`
}
type CreateWalletResponseData struct {
	Wallet *CreatedWallet `json:"wallet"`
}

type CreateWalletResponseBody = Response[CreateWalletResponseData]

// CreateWallet Create godoc
// @Summary      Create a new wallet for user.
// @Description  Create a new wallet for user.
// @Tags         wallet
// @Accept       application/json
// @Produce      application/json
// @Param        request body CreateWalletRequestBody true "Create Wallet Request Body"
// @Success      200  {object}  CreateWalletResponseBody
// @Failure      500  {object}  ErrorResponseBody
// @Router       /wallet [post]
func (h Handlers) CreateWallet(w http.ResponseWriter, r *http.Request) {

	form := &CreateWalletRequestBody{}
	json.NewDecoder(r.Body).Decode(form)
	if form.UserName == "" {
		response_types.ErrorNoBody(w, http.StatusBadRequest, errors.New("user name is required"))
	}

	wallet, err := h.service.CreateWallet(r.Context(), form.UserName, form.Currency)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	response_types.OkJsonBody(w, CreateWalletResponseData{Wallet: &CreatedWallet{
		Id:            wallet.Id,
		UserAccountId: wallet.UserAccountId,
		Balance:       wallet.Balance.String(),
		Currency:      wallet.Currency,
	}})
}

type DepositRequestBody struct {
	Amount string `json:"amount" example:"10.23"`
	Nonce  int64  `json:"nonce" example:"1749286345000"`
}

type Ledger struct {
	Id            int64 `json:"id" example:"1"`
	WalletId      int64 `json:"wallet_id" example:"1"`
	TransactionId int64 `json:"transaction_id" example:"1749286345000"`
	//Operation string    `json:"operation" example:"deposit,withdraw,transfer"`
	EntryType string    `json:"entry_type" example:"credit,debit"`
	Amount    string    `json:"amount" example:"40.22"`
	CreatedAt time.Time `json:"created_at"`
	Balance   string    `json:"balance" example:"2.234"`
}

type Transaction struct {
	Ledgers []Ledger `json:"ledgers"`

	Id            int64     `json:"id" example:"1"`
	UserAccountId int64     `json:"user_account_id" example:"1"`
	Nonce         int64     `json:"nonce"`
	Status        string    `json:"status"`
	Operation     string    `json:"operation"`
	CreatedAt     time.Time `json:"created_at"`
}

type DepositResponseData struct {
	Transaction `json:"transaction"`
}

type DepositResponseBody = Response[DepositResponseData]

// Deposit Create godoc
// @Summary      Deposit to wallet
// @Description  Deposit to wallet
// @Tags         wallet
// @Security     BasicAuth
// @Accept       application/json
// @Produce      application/json
// @Param        wallet_id   					path      string  true  "wallet id"
// @Param        request body DepositRequestBody true "Create Deposit Request Body"
// @Success      200  {object}  DepositResponseBody
// @Failure      500  {object}  ErrorResponseBody
// @Router       /wallet/{wallet_id}/deposit [post]
func (h Handlers) Deposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	basicAuthB64, _ := ctx.Value("BASIC_AUTH").(string)
	principal, err := mustExtractUsernameFromBasicAuthValue(basicAuthB64)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusUnauthorized, err)
		return
	}

	_walletId := r.PathValue("wallet_id")
	walletId, err := strconv.Atoi(_walletId)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	form := &DepositRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	amount, err := decimal.NewFromString(form.Amount)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	transaction, ledger, err := h.service.Deposit(ctx, principal, form.Nonce, int64(walletId), amount)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	response_types.OkJsonBody(w, DepositResponseData{
		Transaction: Transaction{
			Ledgers: []Ledger{
				{
					Id:            ledger.Id,
					WalletId:      ledger.WalletId,
					TransactionId: ledger.TransactionId,
					EntryType:     ledger.EntryType,
					Amount:        ledger.Amount.String(),
					CreatedAt:     ledger.CreatedAt,
					Balance:       ledger.Balance.String(),
				},
			},
			Id:            transaction.Id,
			UserAccountId: transaction.UserAccountId,
			Nonce:         transaction.Nonce,
			Status:        transaction.Status,
			Operation:     transaction.Operation,
			CreatedAt:     transaction.CreatedAt,
		},
	})

}

func mustExtractUsernameFromBasicAuthValue(basicAuthB64 string) (string, error) {
	basicAuth, err := base64.StdEncoding.DecodeString(basicAuthB64)
	if err != nil {
		return "", err
	}
	basicAuthSlice := strings.Split(string(basicAuth), ":")
	if len(basicAuthSlice) != 2 {
		return "", errors.New("invalid_basic_auth")
	}
	principal := basicAuthSlice[0]
	if principal == "" {
		return "", errors.New("invalid principal")
	}
	return principal, nil
}

// Types

type Response[T any] struct {
	Data  T       `json:"data"`
	Error *string `json:"error"`
}

var _ = response_types.Response[int](Response[int]{})

type ErrorResponseBody = response_types.Response[struct{}]
