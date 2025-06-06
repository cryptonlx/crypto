package user

import (
	"context"
	"net/http"
	"strconv"

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

// GetWalletBalance godoc
// @Summary      Get balances of user's wallets.
// @Description  Get balances of user's wallets.
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        user_id   					path      string  true  "username"
// @Success      200  {object}  GetUserWalletBalanceResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /user/{user_id}/balance [get]
func (h Handlers) GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	_user_id := r.PathValue("user_id")

	user_id, err := strconv.Atoi(_user_id)
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
	}

	walletBalances, err := h.service.GetUserWalletBalanceByUserId(context.Background(), int64(user_id))
	if err != nil {
		response_types.ErrorNoBody(w, http.StatusBadRequest, err)
		return
	}
	response_types.OkJsonBody(w, GetUserBalanceResponseData{walletBalances: walletBalances})
}

type GetUserWalletBalanceResponse = Response[GetUserBalanceResponseData]
type Response[T any] struct {
	Data  T       `json:"data"`
	Error *string `json:"error"`
}

var _response = response_types.Response[int](Response[int]{})

type ErrorResponse = response_types.Response[struct{}]
