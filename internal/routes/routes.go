/**

 filename  : routes.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-27

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package routes

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"

	"github.com/zuudevs/service-order-api/internal/handlers"
)

func RegisterPersonRoutes(
	r chi.Router,
	personHandler *handlers.PersonHandler,
) {
	r.Route("/persons", func(r chi.Router) {
		r.Post("/", personHandler.Create)
		r.Get("/", personHandler.Index)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", personHandler.GetByID)
			r.Put("/", personHandler.Replace)
			r.Patch("/", personHandler.Update)
			r.Delete("/", personHandler.Delete)
		})
	})
}

func RegisterContactRoutes(
	r chi.Router,
	contactHandler *handlers.ContactHandler,
) {
	r.Route("/contacts", func(r chi.Router) {
		r.Post("/", contactHandler.Create)
		r.Get("/", contactHandler.Index)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", contactHandler.GetByID)
			r.Put("/", contactHandler.Replace)
			r.Patch("/", contactHandler.Update)
			r.Delete("/", contactHandler.Delete)
		})
	})
}

func RegisterOrderRoutes(
	r chi.Router,
	orderHandler *handlers.OrderHandler,
) {
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.Create)
		r.Get("/", orderHandler.Index)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", orderHandler.GetByID)
			r.Put("/", orderHandler.Replace)
			r.Patch("/", orderHandler.Update)
			r.Delete("/", orderHandler.Delete)
		})
	})
}

func RegisterTaskRoutes(
	r chi.Router,
	taskHandler *handlers.TaskHandler,
) {
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", taskHandler.Create)
		r.Get("/", taskHandler.Index)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", taskHandler.GetByID)
			r.Put("/", taskHandler.Replace)
			r.Patch("/", taskHandler.Update)
			r.Delete("/", taskHandler.Delete)
		})
	})
}

func RegisterTransactionRoutes(
	r chi.Router,
	transactionHandler *handlers.TransactionHandler,
) {
	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", transactionHandler.Create)
		r.Get("/", transactionHandler.Index)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", transactionHandler.GetByID)
			r.Put("/", transactionHandler.Replace)
			r.Patch("/", transactionHandler.Update)
			r.Delete("/", transactionHandler.Delete)
		})
	})
}

func RegisterDetailTaskRoutes(
	r chi.Router,
	detailTaskHandler *handlers.DetailTaskHandler,
) {
	r.Route("/detail-tasks", func(r chi.Router) {
		r.Post("/", detailTaskHandler.Create)
		r.Get("/", detailTaskHandler.Index)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", detailTaskHandler.GetByID)
			r.Put("/", detailTaskHandler.Replace)
			r.Patch("/", detailTaskHandler.Update)
			r.Delete("/", detailTaskHandler.Delete)
		})
	})
}

func RegisterHealthRoutes(
	r chi.Router,
) {
	r.Get("/healthz", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}
