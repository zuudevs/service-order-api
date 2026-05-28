/**

 filename  : person_handler.go
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

type PersonHandler struct {
	service *services.PersonService
}

func NewPersonHandler(
	service *services.PersonService,
) *PersonHandler {
	return &PersonHandler{
		service: service,
	}
}

type CreateRequestPerson struct {
	FirstName  string  `json:"firstname"`
	MiddleName *string `json:"middlename"`
	LastName   *string `json:"lastname"`
}

type UpdateRequestPerson struct {
	FirstName  string  `json:"firstname"`
	MiddleName *string `json:"middlename"`
	LastName   *string `json:"lastname"`
}

type PatchRequestPerson struct {
	FirstName  *string `json:"firstname,omitempty"`
	MiddleName *string `json:"middlename,omitempty"`
	LastName   *string `json:"lastname,omitempty"`
}

func (h *PersonHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateRequestPerson

	err := json.NewDecoder(
		r.Body,
	).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.Create(
		req.FirstName,
		req.MiddleName,
		req.LastName,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(
		map[string]any{
			"success": true,
		},
	)
}

func (h *PersonHandler) Index(
	w http.ResponseWriter,
	r *http.Request,
) {
	firstName := r.URL.Query().Get("firstname")
	lastName := r.URL.Query().Get("lastname")

	var persons []models.Person
	var err error

	if firstName != "" {
		persons, err = h.service.GetByFirstName(firstName)
	} else if lastName != "" {
		persons, err = h.service.GetByLastName(lastName)
	} else {
		persons, err = h.service.Index()
	}

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(persons)
}

func (h *PersonHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(
			w,
			"invalid person id",
			http.StatusBadRequest,
		)
		return
	}

	person, err := h.service.GetByID(id)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(person)
}

func (h *PersonHandler) Replace(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid person id", http.StatusBadRequest)
		return
	}

	var req UpdateRequestPerson

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	err = h.service.UpdatePerson(id, req.FirstName, req.MiddleName, req.LastName)
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

func (h *PersonHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid person id", http.StatusBadRequest)
		return
	}

	var req PatchRequestPerson

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	if req.FirstName != nil {
		err = h.service.UpdateFirstName(id, *req.FirstName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.MiddleName != nil {
		err = h.service.UpdateMiddleName(id, *req.MiddleName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.LastName != nil {
		err = h.service.UpdateLastName(id, *req.LastName)
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

func (h *PersonHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid person id", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}

	err = h.service.DeletePerson(id)
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
