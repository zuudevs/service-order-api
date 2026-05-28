/**

 filename  : detail_task_handler.go
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

type DetailTaskHandler struct {
	service *services.DetailTaskService
}

func NewDetailTaskHandler(
	service *services.DetailTaskService,
) *DetailTaskHandler {
	return &DetailTaskHandler{
		service: service,
	}
}

type CreateDetailTaskRequest struct {
	TaskID uint64 `json:"task_id"`
}

type UpdateDetailTaskRequest struct {
	TaskID uint64 `json:"task_id"`
}

type PatchDetailTaskRequest struct {
	TaskID *uint64 `json:"task_id,omitempty"`
}

func (h *DetailTaskHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateDetailTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = h.service.CreateDetailTask(req.TaskID)

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

func (h *DetailTaskHandler) Index(
	w http.ResponseWriter,
	r *http.Request,
) {
	taskIDStr := r.URL.Query().Get("task_id")

	var detailTasks []models.DetailTask
	var err error

	if taskIDStr != "" {
		taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid task_id", http.StatusBadRequest)
			return
		}
		detailTasks, err = h.service.GetByTaskID(taskID)
	} else {
		detailTasks, err = h.service.Index()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(detailTasks)
}

func (h *DetailTaskHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid detail task id", http.StatusBadRequest)
		return
	}

	detailTask, err := h.service.GetByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(detailTask)
}

func (h *DetailTaskHandler) Replace(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid detail task id", http.StatusBadRequest)
		return
	}

	var req UpdateDetailTaskRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "detail task not found", http.StatusNotFound)
		return
	}

	err = h.service.UpdateTaskID(id, req.TaskID)
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

func (h *DetailTaskHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid detail task id", http.StatusBadRequest)
		return
	}

	var req PatchDetailTaskRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "detail task not found", http.StatusNotFound)
		return
	}

	if req.TaskID != nil {
		err = h.service.UpdateTaskID(id, *req.TaskID)
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

func (h *DetailTaskHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid detail task id", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "detail task not found", http.StatusNotFound)
		return
	}

	err = h.service.DeleteDetailTask(id)
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
