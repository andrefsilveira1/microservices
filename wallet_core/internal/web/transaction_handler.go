package web

import (
	"encoding/json"
	"net/http"

	createtransaction "github.com/andrefsilveira1/microservices/internal/usecase/create_transaction"
)

type WebTransactionHandler struct {
	CreateTransacstionUseCase createtransaction.CreateTransactionUseCase
}

func NewWebTransactionHandler(createTransacstionUseCase createtransaction.CreateTransactionUseCase) *WebTransactionHandler {
	return &WebTransactionHandler{
		CreateTransacstionUseCase: createTransacstionUseCase,
	}
}

func (h *WebTransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var dto createtransaction.CreateTransactionInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	output, err := h.CreateTransacstionUseCase.Execute(ctx, dto)
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
