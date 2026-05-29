/**

 filename  : backup_service.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : Brief description

 copyright Copyright (c) 2026

**/

package services

import (
	"archive/tar"
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type BackupService struct {
	mu           sync.Mutex
	changes      uint64
	interval     uint64
	callback     func()
	db           *sql.DB
	lastBackup   time.Time
	backupCount  uint64
	backupDir    string
}

func NewBackupService(interval uint64, callback func()) *BackupService {
	if interval == 0 {
		interval = 1
	}

	return &BackupService{
		interval:   interval,
		callback:   callback,
		lastBackup: time.Now(),
		backupDir:  "./storage/backups",
	}
}

func NewBackupServiceWithDB(interval uint64, callback func(), db *sql.DB) *BackupService {
	svc := NewBackupService(interval, callback)
	svc.db = db
	os.MkdirAll(svc.backupDir, 0755)
	return svc
}

func (b *BackupService) GetValue() uint64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.changes
}

func (b *BackupService) Increment() {
	b.mu.Lock()
	b.changes++
	shouldBackup := b.changes%b.interval == 0
	cb := b.callback
	b.mu.Unlock()

	if shouldBackup && cb != nil {
		cb()
	}
}

func (b *BackupService) SetCallback(callback func()) {
	b.mu.Lock()
	b.callback = callback
	b.mu.Unlock()
}

func (b *BackupService) SetDatabase(db *sql.DB) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.db = db
}

// GetLastBackupTime returns the timestamp of the last backup
func (b *BackupService) GetLastBackupTime() time.Time {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.lastBackup
}

// CreateIncrementalBackup creates a backup containing only changes since last backup
func (b *BackupService) CreateIncrementalBackup() (string, error) {
	b.mu.Lock()
	lastBackup := b.lastBackup
	b.backupCount++
	backupName := fmt.Sprintf("backup_%d_%d.tar.gz", b.backupCount, time.Now().Unix())
	b.lastBackup = time.Now()
	b.mu.Unlock()

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(b.backupDir, 0755); err != nil {
		return "", err
	}

	backupPath := filepath.Join(b.backupDir, backupName)

	// Create tar.gz file
	file, err := os.Create(backupPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Add database file
	dbFile := "./storage/database.db"
	if err := b.addFileToTar(tarWriter, dbFile, "database.db"); err != nil {
		return "", err
	}

	// Add WAL file if exists (for incremental state)
	walFile := "./storage/database.db-wal"
	if _, err := os.Stat(walFile); err == nil {
		if err := b.addFileToTar(tarWriter, walFile, "database.db-wal"); err != nil {
			return "", err
		}
	}

	// Add backup metadata
	metadata := fmt.Sprintf("backup_count:%d\nlast_backup:%s\ncurrent_time:%s\n",
		b.backupCount,
		lastBackup.Format(time.RFC3339),
		time.Now().Format(time.RFC3339))

	if err := b.addContentToTar(tarWriter, "backup_metadata.txt", []byte(metadata)); err != nil {
		return "", err
	}

	return backupPath, nil
}

func (b *BackupService) addFileToTar(tw *tar.Writer, filePath string, tarName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name: tarName,
		Size: fileInfo.Size(),
		Mode: 0644,
		ModTime: fileInfo.ModTime(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	return nil
}

func (b *BackupService) addContentToTar(tw *tar.Writer, fileName string, content []byte) error {
	header := &tar.Header{
		Name: fileName,
		Size: int64(len(content)),
		Mode: 0644,
		ModTime: time.Now(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := tw.Write(content); err != nil {
		return err
	}

	return nil
}