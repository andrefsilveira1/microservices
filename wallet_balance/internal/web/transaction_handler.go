package web

import (
	"encoding/json"
	"net/http"

	gettransaction "github.com/andrefsilveira1/microservices/wallet_balance/internal/usecase/find_transaction"
)

type TransactionHandler struct {
	FindTransactionUseCase gettransaction.FindTransactionUseCase
}

func NewWebTransactionHandler(findtransactionUseCase gettransaction.FindTransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		FindTransactionUseCase: findtransactionUseCase,
	}
}

func (h *TransactionHandler) FindTransaction(w http.ResponseWriter, r *http.Request) {
	var dto gettransaction.FindTransactionInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	output, err := h.FindTransactionUseCase.Execute(ctx, dto)

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
