package services

import (
	"database/sql"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"resizer/models"

	"encoding/base64"

	"github.com/nfnt/resize"
)

type ImageService struct {
	db *sql.DB
}

func NewImageService(db *sql.DB) *ImageService {
	return &ImageService{db: db}
}

func init() {
	// Set Malaysia timezone (UTC+8)
	loc, err := time.LoadLocation("Asia/Kuala_Lumpur")
	if err != nil {
		log.Printf("Error loading timezone: %v, falling back to UTC+8", err)
		// Fallback to fixed UTC+8 offset if timezone data is not available
		loc = time.FixedZone("MYT", 8*60*60)
	}
	time.Local = loc
}

func (i *ImageService) CreateImageTask(projectID int64, imagePath string, targetWidth, targetHeight int, scheduledFor time.Time) (*models.ImageTask, error) {
	task := &models.ImageTask{
		ProjectID:    projectID,
		ImagePath:    imagePath,
		TargetWidth:  targetWidth,
		TargetHeight: targetHeight,
		Status:       "pending",
		CreatedAt:    time.Now(),
		ScheduledFor: scheduledFor,
	}

	result, err := i.db.Exec(`
		INSERT INTO image_tasks (project_id, image_path, target_width, target_height, status, created_at, scheduled_for)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, task.ProjectID, task.ImagePath, task.TargetWidth, task.TargetHeight, task.Status, task.CreatedAt, task.ScheduledFor)

	if err != nil {
		return nil, fmt.Errorf("failed to create image task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get task ID: %w", err)
	}

	task.ID = id
	return task, nil
}

func (i *ImageService) GetPendingTasks() ([]models.ImageTask, error) {
	now := time.Now()

	rows, err := i.db.Query(`
		SELECT id, project_id, image_path, target_width, target_height, status, created_at, scheduled_for
		FROM image_tasks 
		WHERE status = 'pending' 
		AND datetime(scheduled_for, 'localtime') <= datetime(?, 'localtime')
		ORDER BY scheduled_for ASC
	`, now.Format("2006-01-02 15:04:05"))

	if err != nil {
		log.Printf("Error querying pending tasks: %v", err)
		return nil, fmt.Errorf("failed to get pending tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.ImageTask
	for rows.Next() {
		var task models.ImageTask
		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.ImagePath,
			&task.TargetWidth,
			&task.TargetHeight,
			&task.Status,
			&task.CreatedAt,
			&task.ScheduledFor,
		)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if task.ScheduledFor.After(now) {
			log.Printf("Skipping task %d: scheduled for %v MYT, current time %v MYT",
				task.ID,
				task.ScheduledFor.Format("2006-01-02 15:04:05 MST"),
				now.Format("2006-01-02 15:04:05 MST"),
			)
			continue
		}

		log.Printf("Found pending task: ID=%d, ScheduledFor=%v MYT, CurrentTime=%v MYT",
			task.ID,
			task.ScheduledFor.Format("2006-01-02 15:04:05 MST"),
			now.Format("2006-01-02 15:04:05 MST"),
		)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (i *ImageService) ProcessImage(task *models.ImageTask) error {
	log.Printf("Starting to process image task %d", task.ID)

	// Update status to processing
	_, err := i.db.Exec("UPDATE image_tasks SET status = 'processing' WHERE id = ?", task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	log.Printf("Updated task %d status to processing", task.ID)

	log.Printf("Opening source image: %s", task.ImagePath)
	file, err := os.Open(task.ImagePath)
	if err != nil {
		i.updateTaskStatus(task.ID, "failed")
		return fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	// Decode gambar
	var img image.Image
	var decodeErr error

	ext := filepath.Ext(task.ImagePath)
	log.Printf("Decoding image with extension: %s", ext)

	switch ext {
	case ".jpg", ".jpeg":
		img, decodeErr = jpeg.Decode(file)
	case ".png":
		img, decodeErr = png.Decode(file)
	default:
		i.updateTaskStatus(task.ID, "failed")
		return fmt.Errorf("unsupported image format")
	}

	if decodeErr != nil {
		i.updateTaskStatus(task.ID, "failed")
		return fmt.Errorf("failed to decode image: %w", decodeErr)
	}
	log.Printf("Successfully decoded image for task %d", task.ID)

	// Resize gambar
	log.Printf("Resizing image to %dx%d", task.TargetWidth, task.TargetHeight)
	resized := resize.Resize(uint(task.TargetWidth), uint(task.TargetHeight), img, resize.Lanczos3)
	log.Printf("Successfully resized image for task %d", task.ID)

	// Create output directory kalau tak wujud
	outputDir := filepath.Join(filepath.Dir(filepath.Dir(task.ImagePath)), "resized")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		i.updateTaskStatus(task.ID, "failed")
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	log.Printf("Created output directory: %s", outputDir)

	// Create output file
	outputPath := filepath.Join(outputDir, filepath.Base(task.ImagePath))
	log.Printf("Creating output file: %s", outputPath)
	out, err := os.Create(outputPath)
	if err != nil {
		i.updateTaskStatus(task.ID, "failed")
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Save the resized image
	log.Printf("Saving resized image for task %d", task.ID)
	switch filepath.Ext(task.ImagePath) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(out, resized, nil)
	case ".png":
		err = png.Encode(out, resized)
	}

	if err != nil {
		i.updateTaskStatus(task.ID, "failed")
		return fmt.Errorf("failed to save resized image: %w", err)
	}
	log.Printf("Successfully saved resized image for task %d", task.ID)

	// Update status to completed
	return i.updateTaskStatus(task.ID, "completed")
}

func (i *ImageService) updateTaskStatus(taskID int64, status string) error {
	_, err := i.db.Exec("UPDATE image_tasks SET status = ? WHERE id = ?", status, taskID)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}

func (i *ImageService) GetProjectTasks(projectID int64) ([]models.ImageTask, error) {
	rows, err := i.db.Query(`
		SELECT id, project_id, image_path, target_width, target_height, status, created_at, scheduled_for
		FROM image_tasks 
		WHERE project_id = ?
		ORDER BY created_at DESC
	`, projectID)

	if err != nil {
		return nil, fmt.Errorf("failed to get project tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.ImageTask
	for rows.Next() {
		var task models.ImageTask
		err := rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.ImagePath,
			&task.TargetWidth,
			&task.TargetHeight,
			&task.Status,
			&task.CreatedAt,
			&task.ScheduledFor,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (i *ImageService) SaveUploadedFile(projectID int64, fileData []byte, fileName string) (string, error) {

	var projectLocation string
	err := i.db.QueryRow("SELECT location FROM projects WHERE id = ?", projectID).Scan(&projectLocation)
	if err != nil {
		return "", fmt.Errorf("failed to get project location: %w", err)
	}

	// Create uploads folder dalam project location
	uploadsDir := filepath.Join(projectLocation, "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(fileName)
	uniqueFileName := fmt.Sprintf("%d_%d%s", projectID, timestamp, ext)
	filePath := filepath.Join(uploadsDir, uniqueFileName)

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filePath, nil
}

func (i *ImageService) GetImageData(filePath string) (string, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// Convert to base64
	return base64.StdEncoding.EncodeToString(data), nil
}

func (i *ImageService) GetResizedImageData(filePath string) (string, error) {

	outputDir := filepath.Join(filepath.Dir(filepath.Dir(filePath)), "resized")
	resizedPath := filepath.Join(outputDir, filepath.Base(filePath))

	data, err := os.ReadFile(resizedPath)
	if err != nil {
		return "", fmt.Errorf("failed to read resized image file: %w", err)
	}

	// Convert to base64
	return base64.StdEncoding.EncodeToString(data), nil
}
