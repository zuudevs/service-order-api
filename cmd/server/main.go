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

func syncDatabase(
	driveService *gdrive.DriveService,
) {
	filePath := "./storage/database.db"
	fileId := os.Getenv("GOOGLE_DRIVE_DB_FILE_ID")

	// Check if file exists locally
	_, err := os.Stat(filePath)
	fileExists := err == nil

	// If file ID is set, try to download from GDrive
	if fileId != "" {
		log.Println("Downloading database from Google Drive...")
		err := driveService.DownloadFile(fileId, filePath)
		if err != nil {
			log.Printf("Failed to download database: %v. Using existing or creating new database.", err)
		} else {
			log.Println("Database downloaded successfully from Google Drive")
			return
		}
	}

	// If no file exists locally and no file ID, create new and upload
	if !fileExists && fileId == "" {
		log.Println("No database file found. Creating new database...")
		// Database will be created by ConnectSQLite
		return
	}

	// If local file exists but no file ID in env, we'll upload on backup
}

func backup(
	driveService *gdrive.DriveService,
) {
	filePath := "./storage/database.db"

	fileId := os.Getenv("GOOGLE_DRIVE_DB_FILE_ID")

	// ================= Upload First Time =================

	if fileId == "" {

		response, err := driveService.UploadFile(
			filePath,
		)

		if err != nil {
			log.Printf(
				"backup upload failed: %v",
				err,
			)
			return
		}

		log.Println(
			"uploaded new backup:",
			response.Id,
		)

		log.Println(
			"save this ID to GOOGLE_DRIVE_DB_FILE_ID",
		)

		log.Println(
			"GOOGLE_DRIVE_DB_FILE_ID=" + response.Id,
		)

		return
	}

	// ================= Update Existing File =================

	_, err := driveService.UpdateFile(
		fileId,
		filePath,
	)

	if err != nil {
		log.Printf(
			"backup update failed: %v",
			err,
		)
		return
	}

	log.Println(
		"backup updated successfully",
	)
}

func main() {
	// =============================== Google Drive Sync (Download DB if exists) ===============================

	driveSvc, err := gdrive.NewDriveService()

	if err != nil {
		log.Printf(
			"drive init failed: %v (continuing with local database)",
			err,
		)

	} else {
		syncDatabase(driveSvc)
	}

	// =============================== SQLite ===============================

	db, err := database.ConnectSQLite()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// =============================== Google Drive Backup (Periodic Upload) ===============================

	if driveSvc != nil {
		go func() {
			backup(driveSvc)
			ticker := time.NewTicker(6 * time.Hour)
			defer ticker.Stop()

			for range ticker.C {
				backup(driveSvc)
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
	routes.RegisterHealthRoutes(router)

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