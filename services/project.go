package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"resizer/models"
)

type ProjectService struct {
	db *sql.DB
}

func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{db: db}
}

func getDefaultBaseDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create base directory dalam folder Documents
	baseDir := filepath.Join(homeDir, "Documents", "ImageResizer")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create base directory: %w", err)
	}

	return baseDir, nil
}

func (p *ProjectService) CreateProject(name, description string) (*models.Project, error) {
	// Guna project name untuk folder name
	baseDir, err := getDefaultBaseDir()
	if err != nil {
		return nil, err
	}

	safeName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	location := filepath.Join(baseDir, safeName)

	// Create the project directory kalau takda
	if err := os.MkdirAll(location, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	project := &models.Project{
		Name:         name,
		Description:  description,
		CreationTime: time.Now(),
		Location:     location,
	}

	result, err := p.db.Exec(`
		INSERT INTO projects (name, description, creation_time, location)
		VALUES (?, ?, ?, ?)
	`, project.Name, project.Description, project.CreationTime, project.Location)

	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get project ID: %w", err)
	}

	project.ID = id
	return project, nil
}

func (p *ProjectService) GetProject(id int64) (*models.Project, error) {
	project := &models.Project{}
	err := p.db.QueryRow(`
		SELECT id, name, description, creation_time, location
		FROM projects WHERE id = ?
	`, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreationTime,
		&project.Location,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

func (p *ProjectService) ListProjects() ([]models.Project, error) {
	rows, err := p.db.Query(`
		SELECT id, name, description, creation_time, location
		FROM projects ORDER BY creation_time DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreationTime,
			&project.Location,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (p *ProjectService) UpdateProject(project *models.Project) error {
	_, err := p.db.Exec(`
		UPDATE projects 
		SET name = ?, description = ?
		WHERE id = ?
	`, project.Name, project.Description, project.ID)

	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

func (p *ProjectService) DeleteProject(id int64) error {
	// transaction db start
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get project location
	var projectLocation string
	err = tx.QueryRow("SELECT location FROM projects WHERE id = ?", id).Scan(&projectLocation)
	if err != nil {
		return fmt.Errorf("failed to get project location: %w", err)
	}

	// Buang image tasks sebab foreign key constraint
	_, err = tx.Exec("DELETE FROM image_tasks WHERE project_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete project image tasks: %w", err)
	}

	// Delete the project from db
	_, err = tx.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete project from database: %w", err)
	}

	// Buang the whole folder dalam project
	if err := os.RemoveAll(projectLocation); err != nil {
		return fmt.Errorf("failed to delete project directory: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
