/**

 filename  : contact_handler.go
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

type ContactHandler struct {
	service *services.ContactService
}

func NewContactHandler(
	service *services.ContactService,
) *ContactHandler {
	return &ContactHandler{
		service: service,
	}
}

type CreateContactRequest struct {
	Value       string             `json:"value"`
	ContactType models.ContactType `json:"contact_type"`
	IsMain      bool               `json:"is_main"`
	PersonID    uint64             `json:"person_id"`
}

type UpdateContactRequest struct {
	Value       string             `json:"value"`
	ContactType models.ContactType `json:"contact_type"`
	IsMain      bool               `json:"is_main"`
	PersonID    uint64             `json:"person_id"`
}

type PatchContactRequest struct {
	Value       *string             `json:"value,omitempty"`
	ContactType *models.ContactType `json:"contact_type,omitempty"`
	IsMain      *bool               `json:"is_main,omitempty"`
	PersonID    *uint64             `json:"person_id,omitempty"`
}

func (h *ContactHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateContactRequest

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

	err = h.service.CreateContact(
		req.Value,
		req.ContactType,
		req.IsMain,
		req.PersonID,
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

func (h *ContactHandler) Index(
	w http.ResponseWriter,
	r *http.Request,
) {
	contactTypeStr := r.URL.Query().Get("contact_type")
	personIDStr := r.URL.Query().Get("person_id")
	value := r.URL.Query().Get("value")

	var contacts []models.Contact
	var err error

	if contactTypeStr != "" {
		var contactTypeRaw uint64
		contactTypeRaw, err = strconv.ParseUint(contactTypeStr, 10, 8)
		if err != nil {
			http.Error(w, "invalid contact_type", http.StatusBadRequest)
			return
		}
		contacts, err = h.service.GetByContactType(models.ContactType(contactTypeRaw))
	} else if personIDStr != "" {
		var personID uint64
		personID, err = strconv.ParseUint(personIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid person_id", http.StatusBadRequest)
			return
		}
		contacts, err = h.service.GetByPersonID(personID)
	} else if value != "" {
		contacts, err = h.service.GetByValue(value)
	} else {
		contacts, err = h.service.Index()
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

	json.NewEncoder(w).Encode(contacts)
}

func (h *ContactHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(
			w,
			"invalid contact id",
			http.StatusBadRequest,
		)
		return
	}

	contact, err := h.service.GetByID(id)

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

	json.NewEncoder(w).Encode(contact)
}

func (h *ContactHandler) Replace(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid contact id", http.StatusBadRequest)
		return
	}

	var req UpdateContactRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "contact not found", http.StatusNotFound)
		return
	}

	err = h.service.UpdateContact(id, req.Value, req.ContactType, req.IsMain, req.PersonID)
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

func (h *ContactHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid contact id", http.StatusBadRequest)
		return
	}

	var req PatchContactRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "contact not found", http.StatusNotFound)
		return
	}

	if req.Value != nil {
		err = h.service.UpdateValue(id, *req.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.ContactType != nil {
		err = h.service.UpdateContactType(id, *req.ContactType)
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

func (h *ContactHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id_str := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(id_str, 10, 64)

	if err != nil {
		http.Error(w, "invalid contact id", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetByID(id)
	if err != nil {
		http.Error(w, "contact not found", http.StatusNotFound)
		return
	}

	err = h.service.DeleteContact(id)
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