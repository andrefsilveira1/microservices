package web

import (
	"encoding/json"
	"net/http"

	findbalances "github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/find_balances"
)

type BalanceHandler struct {
	FindBalancesUseCase findbalances.FindBalancesUseCase
}

func NewWebBalanceHandler(findBalancesUseCase findbalances.FindBalancesUseCase) *BalanceHandler {
	return &BalanceHandler{
		FindBalancesUseCase: findBalancesUseCase,
	}
}

func (h *BalanceHandler) FindBalance(w http.ResponseWriter, r *http.Request) {
	var dto findbalances.FindBalancesInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	output, err := h.FindBalancesUseCase.Execute(ctx, dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
