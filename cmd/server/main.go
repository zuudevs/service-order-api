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
	
	_, err := os.Stat(filePath)
	fileExists := err == nil

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

	if !fileExists && fileId == "" {
		log.Println("No database file found. Creating new database...")
		return
	}
}

func backupIncremental(
	driveService *gdrive.DriveService,
	backupService *services.BackupService,
) {
	dbFilePath := "./storage/database.db"
	fileId := os.Getenv("GOOGLE_DRIVE_DB_FILE_ID")

	// ================= Upload/Update Main Database File =================
	
	if fileId == "" {
		// First time: upload raw database file
		response, err := driveService.UploadFile(dbFilePath)
		if err != nil {
			log.Printf("failed to upload main database: %v", err)
			return
		}

		log.Printf("main database uploaded successfully: %s", response.Id)
		log.Printf("save this ID to GOOGLE_DRIVE_DB_FILE_ID: %s", response.Id)
		return
	}

	// Update existing main database file (raw, not compressed)
	_, err := driveService.UpdateFile(fileId, dbFilePath)
	if err != nil {
		log.Printf("failed to update main database: %v", err)
		return
	}

	log.Println("main database updated successfully")

	// ================= Create Incremental Backup Archive =================

	backupPath, err := backupService.CreateIncrementalBackup()
	if err != nil {
		log.Printf("failed to create incremental backup: %v", err)
		return
	}

	backupFolderId := os.Getenv("GOOGLE_DRIVE_BACKUP_FOLDER_ID")

	response, err := driveService.UploadIncrementalBackup(backupPath, backupFolderId)
	if err != nil {
		log.Printf("failed to upload incremental backup: %v", err)
		return
	}

	log.Printf("incremental backup archived: %s", response.Id)

	// Keep only last 7 backups to save storage
	if backupFolderId != "" {
		if err := driveService.DeleteOldBackups(backupFolderId, 7); err != nil {
			log.Printf("warning: failed to clean old backups: %v", err)
		}
	}
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

	// =============================== Google Drive Backup (Periodic Incremental Upload) ===============================

	backupSvc := services.NewBackupServiceWithDB(100, nil, db)

	if driveSvc != nil {
		backupSvc.SetCallback(func() {
			backupIncremental(driveSvc, backupSvc)
		})

		go func() {
			// Initial backup after 1 minute
			time.Sleep(1 * time.Minute)
			backupIncremental(driveSvc, backupSvc)

			// Then backup every 6 hours
			ticker := time.NewTicker(6 * time.Hour)
			defer ticker.Stop()

			for range ticker.C {
				backupIncremental(driveSvc, backupSvc)
			}
		}()
	}

	// ================================= Person =================================

	personRepo := repositories.NewPersonRepository(db)
	personService := services.NewPersonServiceWithBackup(personRepo, backupSvc)
	personHandler := handlers.NewPersonHandler(personService)

	// ================================= Contact =================================

	contactRepo := repositories.NewContactRepository(db)
	contactService := services.NewContactServiceWithBackup(contactRepo, personRepo, backupSvc)
	contactHandler := handlers.NewContactHandler(contactService)

	// ================================= Order =================================

	orderRepo := repositories.NewOrderRepository(db)
	orderService := services.NewOrderServiceWithBackup(orderRepo, personRepo, backupSvc)
	orderHandler := handlers.NewOrderHandler(orderService)

	// ================================= Task =================================

	taskRepo := repositories.NewTaskRepository(db)
	taskService := services.NewTaskServiceWithBackup(taskRepo, backupSvc)
	taskHandler := handlers.NewTaskHandler(taskService)

	// ================================= Transaction =================================

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionServiceWithBackup(transactionRepo, orderRepo, backupSvc)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// ================================= DetailTask =================================

	detailTaskRepo := repositories.NewDetailTaskRepository(db)
	detailTaskService := services.NewDetailTaskServiceWithBackup(detailTaskRepo, taskRepo, backupSvc)
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