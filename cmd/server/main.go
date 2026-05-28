/**

 filename  : main.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-26

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/zuudevs/service-order-api/internal/database"
	"github.com/zuudevs/service-order-api/internal/handlers"
	"github.com/zuudevs/service-order-api/internal/repositories"
	"github.com/zuudevs/service-order-api/internal/routes"
	"github.com/zuudevs/service-order-api/internal/services"
	gdrive "github.com/zuudevs/service-order-api/internal/services/gdrive"
)

func main() {
	db, err := database.ConnectSQLite()
	if err != nil {
		log.Fatal(err)
	}

	// =============================== Google Drive ===============================
	driveSvc, err := gdrive.NewDriveService()
	if err != nil {
		log.Printf("drive init failed: %v", err)
	} else {
		go func() {
			backup := func() {
				_, err := driveSvc.UploadFile("storage/database.db")

				if err != nil {
					log.Printf("backup upload failed: %v", err)
				} else {
					log.Println("backup uploaded successfully")
				}
			}

			backup()

			ticker := time.NewTicker(6 * time.Hour)
			defer ticker.Stop()

			for range ticker.C {
				_, err := driveSvc.UploadFile("storage/database.db")
				if err != nil {
					log.Printf("backup upload failed: %v", err)
				} else {
					log.Println("backup uploaded successfully")
				}
			}
		}()
	}

	// ================================= Person =================================

	personRepo := repositories.NewPersonRepository(db)
	personService := services.NewPersonService(personRepo)
	personHandler := handlers.NewPersonHandler(personService)

	// ================================= Contact =================================

	contactRepo := repositories.NewContactRepository(db)
	contactService := services.NewContactService(contactRepo, personRepo)
	contactHandler := handlers.NewContactHandler(contactService)

	// ================================= Order =================================

	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderService(orderRepo, personRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	// ================================= Task =================================

	taskRepo := repositories.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	// ================================= Transaction =================================

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo, orderRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// ================================= DetailTask =================================

	detailTaskRepo := repositories.NewDetailTaskRepository(db)
	detailTaskService := services.NewDetailTaskService(detailTaskRepo, taskRepo)
	detailTaskHandler := handlers.NewDetailTaskHandler(detailTaskService)

	// ================================= Router Setup =================================

	router := chi.NewRouter()

	routes.RegisterPersonRoutes(router, personHandler)
	routes.RegisterContactRoutes(router, contactHandler)
	routes.RegisterOrderRoutes(router, orderHandler)
	routes.RegisterTaskRoutes(router, taskHandler)
	routes.RegisterTransactionRoutes(router, transactionHandler)
	routes.RegisterDetailTaskRoutes(router, detailTaskHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println(
		"server running on :" + port,
	)

	if err := http.ListenAndServe(
		":"+port,
		router,
	); err != nil {
		log.Fatal(err)
	}
}