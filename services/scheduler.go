package services

import (
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	imageService *ImageService
	stopChan     chan struct{}
	wg           sync.WaitGroup
	isRunning    bool
	mutex        sync.Mutex
}

func NewScheduler(imageService *ImageService) *Scheduler {
	return &Scheduler{
		imageService: imageService,
		stopChan:     make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	s.mutex.Lock()
	if s.isRunning {
		s.mutex.Unlock()
		return
	}
	s.isRunning = true
	s.mutex.Unlock()

	log.Println("Scheduler started")
	s.wg.Add(1)
	go s.run()
}

func (s *Scheduler) Stop() {
	s.mutex.Lock()
	if !s.isRunning {
		s.mutex.Unlock()
		return
	}
	s.isRunning = false
	close(s.stopChan)
	s.mutex.Unlock()

	s.wg.Wait()
}

func (s *Scheduler) run() {
	defer s.wg.Done()

	log.Println("Scheduler running, checking tasks every 5 seconds")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			log.Println("Scheduler stopping")
			return
		case <-ticker.C:
			currentTime := time.Now()
			log.Printf("Scheduler checking for pending tasks at %v MYT", currentTime.Format("2006-01-02 15:04:05 MST"))
			s.processPendingTasks()
		}
	}
}

func (s *Scheduler) processPendingTasks() {
	currentTime := time.Now()
	log.Printf("Checking for pending tasks at %v MYT", currentTime.Format("2006-01-02 15:04:05 MST"))

	tasks, err := s.imageService.GetPendingTasks()
	if err != nil {
		log.Printf("Error getting pending tasks: %v", err)
		return
	}

	log.Printf("Found %d pending tasks at %v MYT", len(tasks), currentTime.Format("2006-01-02 15:04:05 MST"))
	for _, task := range tasks {
		log.Printf("Processing task %d (scheduled for %v MYT, current time: %v MYT)",
			task.ID,
			task.ScheduledFor.Format("2006-01-02 15:04:05 MST"),
			currentTime.Format("2006-01-02 15:04:05 MST"),
		)
		if err := s.imageService.ProcessImage(&task); err != nil {
			log.Printf("Error processing image task %d: %v", task.ID, err)
			continue
		}
		log.Printf("Successfully processed task %d", task.ID)
	}
}
