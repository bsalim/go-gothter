package scanner

import (
	"go-gothter/internal/config"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

type LogScanner struct {
	cfg    *config.Config
	logger *logrus.Logger
}

func NewLogScanner(cfg *config.Config, logger *logrus.Logger) *LogScanner {
	return &LogScanner{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *LogScanner) StartMonitoring() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.logger.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Add(s.cfg.LogFiles.AuthLog)
	if err != nil {
		s.logger.Fatalf("Error adding file to watcher: %v", err)
	}
	// err = watcher.Add(s.cfg.LogFiles.NginxLog)
	// if err != nil {
	// 	s.logger.Fatalf("Error adding file to watcher: %v", err)
	// }

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Handle log file events here
			s.logger.Infof("Log event detected: %v", event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			s.logger.Errorf("Error: %v", err)
		}
	}
}

func (s *LogScanner) CheckForPatterns() bool {
	// Implement pattern matching and return if any patterns are detected
	return false
}
