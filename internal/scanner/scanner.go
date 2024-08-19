package scanner

import (
	"go-gothter/internal/config"
	"log"

	"github.com/fsnotify/fsnotify"
)

type LogScanner struct {
	cfg *config.Config
}

func NewLogScanner(cfg *config.Config) *LogScanner {
	return &LogScanner{cfg: cfg}
}

func (s *LogScanner) StartMonitoring() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(s.cfg.LogFiles.AuthLog)
	if err != nil {
		log.Fatalf("Error adding file to watcher: %v", err)
	}
	err = watcher.Add(s.cfg.LogFiles.NginxLog)
	if err != nil {
		log.Fatalf("Error adding file to watcher: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Handle log file events here
			log.Println("Log event detected:", event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

func (s *LogScanner) CheckForPatterns() bool {
	// Implement pattern matching and return if any patterns are detected
	return false
}
