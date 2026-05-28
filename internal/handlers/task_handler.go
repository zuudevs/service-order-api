/**

 filename  : task_handler.go
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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zuudevs/service-order-api/internal/models"
	"github.com/zuudevs/service-order-api/internal/services"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(
	service *services.TaskService,
) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

type CreateTaskRequest struct {
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Price       uint64    `json:"price"`
	Due         time.Time `json:"due"`
}

type UpdateTaskRequest struct {
	Subject     string            `json:"subject"`
	Description string            `json:"description"`
	Price       uint64            `json:"price"`
	Due         time.Time         `json:"due"`
	Status      models.TaskStatus `json:"status"`
}

type PatchTaskRequest struct {
	Subject     *string            `json:"subject,omitempty"`
	Description *string            `json:"description,omitempty"`
	Price       *uint64            `json:"price,omitempty"`
	Status      *models.TaskStatus `json:"status,omitempty"`
}

func (h *TaskHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	err = h.service.CreateTask(req.Subject, req.Description, req.Price, req.Due)

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

func (h *TaskHandler) Index(
	w http.ResponseWriter,
	r *http.Request,
) {
	statusStr := r.URL.Query().Get("status")
	subject := r.URL.Query().Get("subject")

	var tasks []models.Task
	var err error

	if statusStr != "" {
		status, err := strconv.ParseUint(statusStr, 10, 8)
		if err != nil {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		tasks, err = h.service.GetByStatus(models.TaskStatus(status))
	} else if subject != "" {
		tasks, err = h.service.GetBySubject(subject)
	} else {
		tasks, err = h.service.Index()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Replace(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	err = h.service.UpdateSubject(id, req.Subject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdatePrice(id, req.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.UpdateStatus(id, req.Status)
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

func (h *TaskHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var req PatchTaskRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if req.Subject != nil {
		err = h.service.UpdateSubject(id, *req.Subject)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.Price != nil {
		err = h.service.UpdatePrice(id, *req.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.Status != nil {
		err = h.service.UpdateStatus(id, *req.Status)
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

func (h *TaskHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	err = h.service.DeleteTask(id)
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
