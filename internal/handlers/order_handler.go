/**

 filename  : order_handler.go
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

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(
	service *services.OrderService,
) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

type CreateOrderRequest struct {
	Status   models.OrderStatus `json:"status"`
	PersonID *uint64            `json:"person_id"`
}

type UpdateOrderRequest struct {
	Status     models.OrderStatus `json:"status"`
	TotalPrice uint64             `json:"total_price"`
	PersonID   *uint64            `json:"person_id"`
}

type PatchOrderRequest struct {
	Status     *models.OrderStatus `json:"status,omitempty"`
	TotalPrice *uint64             `json:"total_price,omitempty"`
	PersonID   *uint64             `json:"person_id,omitempty"`
}

func (h *OrderHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateOrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = h.service.CreateOrder(req.Status, req.PersonID)

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

func (h *OrderHandler) Index(
	w http.ResponseWriter,
	r *http.Request,
) {
	statusStr := r.URL.Query().Get("status")
	personIDStr := r.URL.Query().Get("person_id")

	var orders []models.Order
	var err error

	if statusStr != "" {
		var statusRaw uint64
		statusRaw, err = strconv.ParseUint(statusStr, 10, 8)
		if err != nil {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		orders, err = h.service.GetByStatus(models.OrderStatus(statusRaw))
	} else if personIDStr != "" {
		var personID uint64
		personID, err = strconv.ParseUint(personIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid person_id", http.StatusBadRequest)
			return
		}
		orders, err = h.service.GetByPersonID(personID)
	} else {
		orders, err = h.service.Index()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) Replace(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	var req UpdateOrderRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	err = h.service.UpdateStatus(id, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateTotalPrice(id, req.TotalPrice)
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

func (h *OrderHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	var req PatchOrderRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	if req.Status != nil {
		err = h.service.UpdateStatus(id, *req.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.TotalPrice != nil {
		err = h.service.UpdateTotalPrice(id, *req.TotalPrice)
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

func (h *OrderHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	err = h.service.DeleteOrder(id)
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