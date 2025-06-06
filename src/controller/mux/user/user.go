package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cryptonlx/crypto/src/controller/response_types"
	userrepo "github.com/cryptonlx/crypto/src/repositories/user"
	userservice "github.com/cryptonlx/crypto/src/service/user"
)

type Handlers struct {
	service *userservice.Service
}

func NewHandlers(service *userservice.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

type GetUserBalanceResponseData struct {
	walletBalances userrepo.UserBalance `json:"wallet_balances"`
}

type User struct {
	Id       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"user1"`
}

// GetWalletBalance godoc
// @Summary      Get balances of user's wallets.
// @Description  Get balances of user's wallets.
// @Tags         wallet
// @Accept       application/json
// @Produce      application/json
// @Param        user_id   					path      string  true  "username"
// @Success      200  {object}  GetUserWalletBalanceResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /user/{username}/balance [get]
func (h Handlers) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("username")

	walletBalances, err := h.service.GetUserWalletBalanceByUserName(r.Context(), userName)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	response_types.OkJsonBody(w, GetUserBalanceResponseData{walletBalances: walletBalances})
}

type CreateRequestBody struct {
	UserName string `json:"username"` // Subreddit Name
}

type CreateUserResponseData struct {
	User User `json:"user"`
}

// CreateUser Create godoc
// @Summary      Create a new user.
// @Description  Create a new user.
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        request body CreateRequestBody true "Create Request Body"
// @Success      200  {object}  CreateUserResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /user [post]
func (h Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	form := &CreateRequestBody{}
	json.NewDecoder(r.Body).Decode(form)
	if form.UserName == "" {
		response_types.ErrorNoBody(w, http.StatusBadRequest, errors.New("user name is required"))
	}

	user, err := h.service.CreateUser(r.Context(), form.UserName)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}

	response_types.OkJsonBody(w, CreateUserResponseData{User: User(user)})
}

type CreateUserResponse = Response[CreateUserResponseData]
type GetUserWalletBalanceResponse = Response[GetUserBalanceResponseData]
type Response[T any] struct {
	Data  T       `json:"data"`
	Error *string `json:"error"`
}

var _ = response_types.Response[int](Response[int]{})

type ErrorResponse = response_types.Response[struct{}]
