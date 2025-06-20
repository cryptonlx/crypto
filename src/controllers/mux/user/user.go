package user

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cryptonlx/crypto/src/controllers/response_types"
	userservice "github.com/cryptonlx/crypto/src/services/user"

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

type GetWalletsResponseBody = ResponseBody[GetWalletsResponseData]

// Wallets godoc
// @Summary      Get balances of user's wallets.
// @Description  Get balances of user's wallets.
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        user_id   					path      string  true  "username"
// @Success      200  {object}  GetWalletsResponseBody
// @Failure      500  {object}  ErrorResponseBody500
// @Router       /user/{username}/wallets [get]
func (h Handlers) Wallets(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("username")

	walletBalances, err := h.service.GetUserWalletBalanceByUserName(r.Context(), userName)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
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

	response_types.WriteOkJsonBody(w, GetWalletsResponseData{
		User:    User(walletBalances.User),
		Wallets: wallets,
	})
}

type CreateUserRequestBody struct {
	UserName string `json:"username" example:"user1"`
}

type CreatedUser struct {
	Id       int64  `json:"id" example:"102"`
	Username string `json:"username" example:"user1"`
}
type CreateUserResponseData struct {
	User *CreatedUser `json:"user"`
}

type CreateUserResponseBody = ResponseBody[CreateUserResponseData]

// CreateUser Create godoc
// @Summary      Create a new user.
// @Description  Create a new user.
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        request body CreateUserRequestBody true "Create User Request Body"
// @Success      200  {object}  CreateUserResponseBody
// @Failure      400  {object}  ErrorResponseBody400
// @Failure      500  {object}  ErrorResponseBody500
// @Router       /user [post]
func (h Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	form := &CreateUserRequestBody{}
	json.NewDecoder(r.Body).Decode(form)
	if form.UserName == "" {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, errors.New("user name is required"))
	}

	user, err := h.service.CreateUser(r.Context(), form.UserName)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	c := CreatedUser(user)
	response_types.WriteOkJsonBody(w, CreateUserResponseData{User: &c})
}

type TransactionMetaData struct {
	SourceWalletId *int64  `json:"source_wallet_id" example:"1021"`
	Amount         *string `json:"amount" example:"40.1122"`
}

type Transaction struct {
	Ledgers []Ledger `json:"ledgers"`

	Id          int64     `json:"id" example:"1"`
	RequestorId int64     `json:"requestor_id" example:"1"`
	Nonce       int64     `json:"nonce" example:"1749460653395"`
	Status      string    `json:"status" example:"success"`
	Operation   string    `json:"operation" example:"deposit"`
	CreatedAt   time.Time `json:"created_at" example:"2025-06-09T02:02:31.213543+08:00"`

	TransactionMetaData `json:"metadata"`
}

type UserTransactions struct {
	User         User
	Transactions []Transaction
}

type TransactionsResponseData struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionsResponseBody = ResponseBody[TransactionsResponseData]

// Transactions godoc
// @Summary      Get transactions of user's wallets sorted by newest.
// @Description  Get transactions of user's wallets sorted by newest.
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        user_id   					path      string  true  "username"
// @Success      200  {object}  TransactionsResponseBody
// @Failure      400  {object}  ErrorResponseBody400
// @Router       /user/{username}/transactions [get]
func (h Handlers) Transactions(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("username")

	transactionLedgers, err := h.service.GetUserTransactionsByUserName(r.Context(), userName)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
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

		t := transaction.Transaction
		var amount *string
		if t.MetaData.Amount != nil {
			amountP := t.MetaData.Amount.String()
			amount = &amountP
		}
		transactions = append(transactions, Transaction{
			Ledgers:     ledgers,
			Id:          t.Id,
			RequestorId: t.RequestorId,
			Nonce:       t.Nonce,
			Status:      t.Status,
			Operation:   t.Operation,
			CreatedAt:   t.CreatedAt,
			TransactionMetaData: TransactionMetaData{
				SourceWalletId: t.MetaData.SourceWalletId,
				Amount:         amount,
			},
		})
	}
	response_types.WriteOkJsonBody(w, TransactionsResponseData{
		Transactions: transactions,
	})
}

type CreateWalletRequestBody struct {
	UserName string `json:"username" example:"username1"`
	Currency string `json:"currency" example:"USD"`
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

type CreateWalletResponseBody = ResponseBody[CreateWalletResponseData]

// CreateWallet Create godoc
// @Summary      Create a new wallet for user.
// @Description  Create a new wallet for user.
// @Tags         wallet
// @Accept       application/json
// @Produce      application/json
// @Param        request body CreateWalletRequestBody true "Create Wallet Request Body"
// @Success      200  {object}  CreateWalletResponseBody
// @Failure      400  {object}  ErrorResponseBody400
// @Failure      500  {object}  ErrorResponseBody500
// @Router       /wallet [post]
func (h Handlers) CreateWallet(w http.ResponseWriter, r *http.Request) {
	form := &CreateWalletRequestBody{}
	json.NewDecoder(r.Body).Decode(form)
	if form.UserName == "" {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, errors.New("user name is required"))
	}

	wallet, err := h.service.CreateWallet(r.Context(), form.UserName, form.Currency)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	response_types.WriteOkJsonBody(w, CreateWalletResponseData{Wallet: &CreatedWallet{
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
	Id            int64     `json:"id" example:"12222214214"`
	WalletId      int64     `json:"wallet_id" example:"1021"`
	TransactionId int64     `json:"transaction_id" example:"1749286345000"`
	EntryType     string    `json:"entry_type" example:"credit"`
	Amount        string    `json:"amount" example:"40.1122"`
	CreatedAt     time.Time `json:"created_at" example:"2025-06-09T02:02:31.213543+08:00"`
	Balance       string    `json:"balance" example:"2.2324"`
}

type DepositResponseData struct {
	Transaction `json:"transaction"`
}

type DepositResponseBody = ResponseBody[DepositResponseData]

// Deposit Create godoc
// @Summary      Deposit to wallet
// @Description  Deposit to wallet
// @Tags         wallet
// @Security     BasicAuth
// @Accept       application/json
// @Produce      application/json
// @Param 		 Authorization header string true "Basic Authorization"
// @Param        wallet_id   					path      string  true  "Wallet Id"
// @Param        request body DepositRequestBody true "Create Deposit Request Body"
// @Success      200  {object}  DepositResponseBody
// @Failure      400  {object}  ErrorResponseBody400
// @Failure      500  {object}  ErrorResponseBody500
// @Router       /wallet/{wallet_id}/deposit [post]
func (h Handlers) Deposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	basicAuthB64, _ := ctx.Value("BASIC_AUTH").(string)
	principal, err := mustExtractUsernameFromBasicAuthValue(basicAuthB64)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusUnauthorized, err)
		return
	}

	_walletId := r.PathValue("wallet_id")
	walletId, err := strconv.Atoi(_walletId)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	form := &DepositRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	amount, err := decimal.NewFromString(form.Amount)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	transaction, ledger, err := h.service.Deposit(ctx, principal, form.Nonce, int64(walletId), amount)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	response_types.WriteOkJsonBody(w, DepositResponseData{
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
			Id:          transaction.Id,
			RequestorId: transaction.RequestorId,
			Nonce:       transaction.Nonce,
			Status:      transaction.Status,
			Operation:   transaction.Operation,
			CreatedAt:   transaction.CreatedAt,
		},
	})
}

type WithdrawRequestBody struct {
	Amount string `json:"amount" example:"10.23"`
	Nonce  int64  `json:"nonce" example:"1749286345000"`
}

var _ = WithdrawLedger(Ledger{})

type WithdrawLedger struct {
	Id            int64     `json:"id" example:"12222214214"`
	WalletId      int64     `json:"wallet_id" example:"1021"`
	TransactionId int64     `json:"transaction_id" example:"1749286345000"`
	EntryType     string    `json:"entry_type" example:"debit"`
	Amount        string    `json:"amount" example:"40.1122"`
	CreatedAt     time.Time `json:"created_at" example:"2025-06-09T02:02:31.213543+08:00"`
	Balance       string    `json:"balance" example:"2.2324"`
}

type WithdrawTransaction struct {
	Ledgers []WithdrawLedger `json:"ledgers"`

	Id          int64     `json:"id" example:"1"`
	RequestorId int64     `json:"requestor_id" example:"1"`
	Nonce       int64     `json:"nonce" example:"1749460653395"`
	Status      string    `json:"status" example:"success"`
	Operation   string    `json:"operation" example:"withdraw"`
	CreatedAt   time.Time `json:"created_at" example:"2025-06-09T02:02:31.213543+08:00"`

	TransactionMetaData `json:"metadata"`
}

type WithdrawResponseData struct {
	WithdrawTransaction `json:"transaction"`
}

type WithdrawResponseBody = ResponseBody[WithdrawResponseData]

// Withdraw Create godoc
// @Summary      Withdraw from wallet
// @Description  Withdraw from wallet
// @Tags         wallet
// @Security     BasicAuth
// @Accept       application/json
// @Produce      application/json
// @Param 		 Authorization header string true "Basic Authorization"
// @Param        wallet_id   					path      string  true  "Wallet Id"
// @Param        request body WithdrawRequestBody true "Create Withdraw Request Body"
// @Success      200  {object}  WithdrawResponseBody
// @Failure      400  {object}  ErrorResponseBody400
// @Failure      500  {object}  ErrorResponseBody500
// @Router       /wallet/{wallet_id}/withdrawal [post]
func (h Handlers) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	basicAuthB64, _ := ctx.Value("BASIC_AUTH").(string)
	principal, err := mustExtractUsernameFromBasicAuthValue(basicAuthB64)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusUnauthorized, err)
		return
	}

	_walletId := r.PathValue("wallet_id")
	walletId, err := strconv.Atoi(_walletId)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	form := &WithdrawRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	amount, err := decimal.NewFromString(form.Amount)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	transaction, ledger, err := h.service.Withdraw(ctx, principal, form.Nonce, int64(walletId), amount)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	var amountP *string
	if transaction.MetaData.Amount != nil {
		_amount := transaction.MetaData.Amount.String()
		amountP = &_amount
	}

	response_types.WriteOkJsonBody(w, DepositResponseData{
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
			Id:          transaction.Id,
			RequestorId: transaction.RequestorId,
			Nonce:       transaction.Nonce,
			Status:      transaction.Status,
			Operation:   transaction.Operation,
			CreatedAt:   transaction.CreatedAt,
			TransactionMetaData: TransactionMetaData{
				SourceWalletId: transaction.MetaData.SourceWalletId,
				Amount:         amountP,
			},
		},
	})
}

type TransferRequestBody struct {
	Amount              string `json:"amount" example:"10.23"`
	Nonce               int64  `json:"nonce" example:"1749286345000"`
	DestinationWalletId int64  `json:"destination_wallet_id" example:"2"`
}

type TransferResponseData struct {
	Transaction `json:"transaction"`
}

type TransferResponseBody = ResponseBody[TransferResponseData]

// Transfer Create godoc
// @Summary      Transfer to another wallet.
// @Description  Transfer to another wallet.
// @Tags         wallet
// @Security     BasicAuth
// @Accept       application/json
// @Produce      application/json
// @Param 		 Authorization header string true "Basic Authorization"
// @Param        wallet_id   					path      string  true  "Wallet Id"
// @Param        request body TransferRequestBody true "Create Transfer Request Body"
// @Success      200  {object}  TransferResponseBody
// @Failure      400  {object}  ErrorResponseBody400
// @Failure      500  {object}  ErrorResponseBody500
// @Router       /wallet/{wallet_id}/transfer [post]
func (h Handlers) Transfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	basicAuthB64, _ := ctx.Value("BASIC_AUTH").(string)
	principal, err := mustExtractUsernameFromBasicAuthValue(basicAuthB64)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusUnauthorized, err)
		return
	}

	_walletId := r.PathValue("wallet_id")
	walletId, err := strconv.Atoi(_walletId)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	form := &TransferRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	amount, err := decimal.NewFromString(form.Amount)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	transaction, ledgersS, err := h.service.Transfer(ctx, principal, form.Nonce, int64(walletId), form.DestinationWalletId, amount)
	if err != nil {
		response_types.WriteErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	var amountP *string
	if transaction.MetaData.Amount != nil {
		_amount := transaction.MetaData.Amount.String()
		amountP = &_amount
	}

	ledgers := make([]Ledger, 0, len(ledgersS))
	for _, ledger := range ledgersS {
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

	response_types.WriteOkJsonBody(w, TransferResponseData{
		Transaction: Transaction{
			Ledgers:     ledgers,
			Id:          transaction.Id,
			RequestorId: transaction.RequestorId,
			Nonce:       transaction.Nonce,
			Status:      transaction.Status,
			Operation:   transaction.Operation,
			CreatedAt:   transaction.CreatedAt,
			TransactionMetaData: TransactionMetaData{
				SourceWalletId: transaction.MetaData.SourceWalletId,
				Amount:         amountP,
			},
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

var _ = response_types.ResponseBody[struct{}](ResponseBody[struct{}]{})

type ResponseBody[T any] struct {
	Data  T       `json:"data"`
	Error *string `json:"error" example:"" extensions:"x-nullable"`
}

type ErrorResponseBody = struct {
	Data  *int    `json:"data" example:"0"`
	Error *string `json:"error" example:"general_error"`
}

type ErrorResponseBody500 = struct {
	Data  *int    `json:"data" example:"0"`
	Error *string `json:"error" example:"internal_server_error"`
}

type ErrorResponseBody400 = struct {
	Data  *int    `json:"data" example:"0"`
	Error *string `json:"error" example:"error_bad_request"`
}
