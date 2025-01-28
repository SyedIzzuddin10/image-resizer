package main

import (
	"context"
	"fmt"
	"time"

	"resizer/models"
	"resizer/services"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	authService    *services.AuthService
	projectService *services.ProjectService
	imageService   *services.ImageService
	scheduler      *services.Scheduler
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	db, err := models.InitDB("resizer.db")
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	// Initialize services
	a.authService = services.NewAuthService(db)
	a.projectService = services.NewProjectService(db)
	a.imageService = services.NewImageService(db)
	a.scheduler = services.NewScheduler(a.imageService)

	// Create default admin user kalau tak exist
	err = a.authService.CreateUser("admin", "admin123")
	if err != nil {
		fmt.Printf("Note: Default admin user might already exist: %v\n", err)
	}

	a.scheduler.Start()
}

func (a *App) shutdown(ctx context.Context) {
	a.scheduler.Stop()
}

func (a *App) Login(username, password string) (*services.LoginResponse, error) {
	return a.authService.Login(username, password)
}

func (a *App) CreateProject(name, description string) (*models.Project, error) {
	return a.projectService.CreateProject(name, description)
}

func (a *App) GetProject(id int64) (*models.Project, error) {
	return a.projectService.GetProject(id)
}

func (a *App) ListProjects() ([]models.Project, error) {
	return a.projectService.ListProjects()
}

func (a *App) UpdateProject(project *models.Project) error {
	return a.projectService.UpdateProject(project)
}

func (a *App) DeleteProject(id int64) error {
	return a.projectService.DeleteProject(id)
}

func (a *App) CreateImageTask(projectID int64, imagePath string, targetWidth, targetHeight int, scheduledFor string) (*models.ImageTask, error) {
	scheduledTime, err := time.Parse(time.RFC3339, scheduledFor)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduled time format: %w", err)
	}
	return a.imageService.CreateImageTask(projectID, imagePath, targetWidth, targetHeight, scheduledTime)
}

func (a *App) GetProjectTasks(projectID int64) ([]models.ImageTask, error) {
	return a.imageService.GetProjectTasks(projectID)
}

func (a *App) SaveUploadedFile(projectID int64, fileData []byte, fileName string) (string, error) {
	return a.imageService.SaveUploadedFile(projectID, fileData, fileName)
}

func (a *App) GetImageData(filePath string) (string, error) {
	return a.imageService.GetImageData(filePath)
}

func (a *App) GetResizedImageData(filePath string) (string, error) {
	return a.imageService.GetResizedImageData(filePath)
}

func (a *App) MessageDialog(title string, message string, dialogType string) bool {
	options := runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         title,
		Message:       message,
		Buttons:       []string{"Cancel", "OK"},
		DefaultButton: "Cancel",
		CancelButton:  "Cancel",
	}

	switch dialogType {
	case "info":
		options.Type = runtime.InfoDialog
	case "warning":
		options.Type = runtime.WarningDialog
	case "error":
		options.Type = runtime.ErrorDialog
	case "question":
		options.Type = runtime.QuestionDialog
	}

	result, err := runtime.MessageDialog(a.ctx, options)
	if err != nil {
		fmt.Printf("Error showing dialog: %v\n", err)
		return false
	}

	return result == "OK"
}
