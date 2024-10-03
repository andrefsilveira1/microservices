package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	findbalances "github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/find_balances"
	"github.com/gorilla/mux"
)

type BalanceHandler struct {
	FindBalancesUseCase findbalances.FindBalancesUseCase
}

func NewWebBalanceHandler(findBalancesUseCase findbalances.FindBalancesUseCase) *BalanceHandler {
	return &BalanceHandler{
		FindBalancesUseCase: findBalancesUseCase,
	}
}

func (h *BalanceHandler) FindBalanceByID(id findbalances.FindBalancesInputDTO, ctx context.Context) (*findbalances.FindBalancesOutputDTO, error) {
	// Call the use case to get the balance by id
	return h.FindBalancesUseCase.Execute(ctx, id)
}

func (h *BalanceHandler) FindBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("GO ID INPUT ===> ", id)
	dto := findbalances.FindBalancesInputDTO{
		ID: id,
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
