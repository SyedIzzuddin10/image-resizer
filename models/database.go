package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Project struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreationTime time.Time `json:"creation_time"`
	Location     string    `json:"location"`
}

type ImageTask struct {
	ID           int64     `json:"id"`
	ProjectID    int64     `json:"project_id"`
	ImagePath    string    `json:"image_path"`
	TargetWidth  int       `json:"target_width"`
	TargetHeight int       `json:"target_height"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	ScheduledFor time.Time `json:"scheduled_for"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password_hash"`
}

func InitDB(dbPath string) (*sql.DB, error) {
	// Create the database file if it doesn't exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.OpenFile(dbPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %w", err)
		}
		file.Close()

		// Ensure the file has the correct permissions
		if err := os.Chmod(dbPath, 0666); err != nil {
			return nil, fmt.Errorf("failed to set database permissions: %w", err)
		}
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set busy timeout to handle concurrent access
	_, err = db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	// Create projects table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			creation_time DATETIME NOT NULL,
			location TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create projects table: %w", err)
	}

	// Create image_tasks table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS image_tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			image_path TEXT NOT NULL,
			target_width INTEGER NOT NULL,
			target_height INTEGER NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			scheduled_for DATETIME,
			FOREIGN KEY (project_id) REFERENCES projects (id)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create image_tasks table: %w", err)
	}

	return db, nil
}
