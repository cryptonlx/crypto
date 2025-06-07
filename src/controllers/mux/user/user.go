package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cryptonlx/crypto/src/controllers/response_types"
	userrepo "github.com/cryptonlx/crypto/src/repositories/user"
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

type GetWalletBalanceResponseData struct {
	User    User     `json:"user"`
	Wallets []Wallet `json:"wallets"`
}

type GetUserWalletBalanceResponse = Response[GetWalletBalanceResponseData]

// GetWalletBalances godoc
// @Summary      Get balances of user's wallets.
// @Description  Get balances of user's wallets.
// @Tags         wallet
// @Accept       application/json
// @Produce      application/json
// @Param        user_id   					path      string  true  "username"
// @Success      200  {object}  GetUserWalletBalanceResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /user/{username}/balance [get]
func (h Handlers) GetWalletBalances(w http.ResponseWriter, r *http.Request) {
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

	response_types.OkJsonBody(w, GetWalletBalanceResponseData{
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
// @Failure      500  {object}  ErrorResponse
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

type GetUserTransactionsResponseData struct {
	transactions userrepo.UserTransactions `json:"wallet_balances"`
}

type GetUserTransactionsResponse = Response[GetUserTransactionsResponseData]

// GetWalletTransactions godoc
// @Summary      Get transactions of user's wallets.
// @Description  Get transactions of user's wallets.
// @Tags         wallet
// @Accept       application/json
// @Produce      application/json
// @Param        user_id   					path      string  true  "username"
// @Success      200  {object}  GetUserTransactionsResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /user/{username}/transactions [get]
func (h Handlers) GetWalletTransactions(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("username")

	transactions, err := h.service.GetUserTransactionsByUserName(r.Context(), userName)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	response_types.OkJsonBody(w, GetUserTransactionsResponseData{
		transactions: transactions,
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
// @Failure      500  {object}  ErrorResponse
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

type Deposit struct {
	Id int64 `json:"id" example:"1"`
}
type DepositResponseData struct {
	Ledger userrepo.Ledger `json:"ledger"`
}

type DepositResponseBody = Response[DepositResponseData]

// Deposit Create godoc
// @Summary      Deposit to wallet
// @Description  Deposit to wallet
// @Tags         wallet
// @Accept       application/json
// @Produce      application/json
// @Param        wallet_id   					path      string  true  "wallet id"
// @Param        request body DepositRequestBody true "Create Deposit Request Body"
// @Success      200  {object}  DepositResponseBody
// @Failure      500  {object}  ErrorResponse
// @Router       /wallet [post]
func (h Handlers) Deposit(w http.ResponseWriter, r *http.Request) {
	_walletId := r.PathValue("wallet_id")
	walletId, err := strconv.Atoi(_walletId)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	form := &DepositRequestBody{}
	json.NewDecoder(r.Body).Decode(form)

	amount, err := decimal.NewFromString(form.Amount)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	ledger, err := h.service.Deposit(r.Context(), form.Nonce, int64(walletId), amount)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	response_types.OkJsonBody(w, DepositResponseData{
		Ledger: ledger,
	})

}

// Types

type Response[T any] struct {
	Data  T       `json:"data"`
	Error *string `json:"error"`
}

var _ = response_types.Response[int](Response[int]{})

type ErrorResponse = response_types.Response[struct{}]
