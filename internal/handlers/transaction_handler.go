/**

 filename  : transaction_handler.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zuudevs/service-order-api/internal/models"
	"github.com/zuudevs/service-order-api/internal/services"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(
	service *services.TransactionService,
) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

type CreateTransactionRequest struct {
	Status       models.TransactionStatus `json:"status"`
	Method       models.TransactionMethod `json:"method"`
	Amount       uint64                   `json:"amount"`
	EvidencePath string                   `json:"evidence_path"`
	OrderID      *uint64                  `json:"order_id"`
}

type UpdateTransactionRequest struct {
	Status       models.TransactionStatus `json:"status"`
	Method       models.TransactionMethod `json:"method"`
	Amount       uint64                   `json:"amount"`
	EvidencePath string                   `json:"evidence_path"`
	OrderID      *uint64                  `json:"order_id"`
}

type PatchTransactionRequest struct {
	Status       *models.TransactionStatus `json:"status,omitempty"`
	Method       *models.TransactionMethod `json:"method,omitempty"`
	Amount       *uint64                   `json:"amount,omitempty"`
	EvidencePath *string                   `json:"evidence_path,omitempty"`
	OrderID      *uint64                   `json:"order_id,omitempty"`
}

func (h *TransactionHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateTransactionRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = h.service.CreateTransaction(req.Status, req.Method, req.Amount, req.EvidencePath, req.OrderID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}

func (h *TransactionHandler) Index(
	w http.ResponseWriter,
	r *http.Request,
) {
	statusStr := r.URL.Query().Get("status")
	methodStr := r.URL.Query().Get("method")
	orderIDStr := r.URL.Query().Get("order_id")

	var transactions []models.Transaction
	var err error

	if statusStr != "" {
		var statusRaw uint64
		statusRaw, err = strconv.ParseUint(statusStr, 10, 8)
		if err != nil {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		transactions, err = h.service.GetByStatus(models.TransactionStatus(statusRaw))
	} else if methodStr != "" {
		var methodRaw uint64
		methodRaw, err = strconv.ParseUint(methodStr, 10, 8)
		if err != nil {
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		}
		transactions, err = h.service.GetByMethod(models.TransactionMethod(methodRaw))
	} else if orderIDStr != "" {
		var orderID uint64
		orderID, err = strconv.ParseUint(orderIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid order_id", http.StatusBadRequest)
			return
		}
		transactions, err = h.service.GetByOrderID(orderID)
	} else {
		transactions, err = h.service.Index()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.GetByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) Replace(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	var req UpdateTransactionRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	err = h.service.UpdateStatus(id, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateAmount(id, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}

func (h *TransactionHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	var req PatchTransactionRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	if req.Status != nil {
		err = h.service.UpdateStatus(id, *req.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.Amount != nil {
		err = h.service.UpdateAmount(id, *req.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}

func (h *TransactionHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	err = h.service.DeleteTransaction(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}