package main

import (
	"bufio"
	"go-gothter/internal/config"
	"go-gothter/internal/notifier"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
	logger.SetOutput(os.Stdout)

	// Load configuration
	cfg, err := config.LoadConfig("configs/default_config.yaml")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}

	// Initialize notifier
	logger.Info("Starting Gothter...")
	not := notifier.NewNotifier(cfg, logger)

	// Monitor log files
	for {
		patternDetected := monitorLogs(cfg, logger)
		if patternDetected {
			not.SendNotification()
		}

		// Sleep for a while before checking again
		time.Sleep(5 * time.Minute)
	}
}

// monitorLogs monitors log files concurrently
func monitorLogs(cfg *config.Config, logger *logrus.Logger) bool {
	var wg sync.WaitGroup
	patternDetected := false
	results := make(chan bool, 2) // Channel to collect results from goroutines

	// Monitor auth.log
	wg.Add(1)
	go func() {
		defer wg.Done()
		if detectPattern(cfg.LogFiles.AuthLog, cfg.Patterns.AuthFail, logger) {
			results <- true
		}
	}()

	// Monitor nginx.log
	wg.Add(1)
	go func() {
		defer wg.Done()
		if detectPattern(cfg.LogFiles.NginxLog, cfg.Patterns.Nginx404, logger) {
			results <- true
		}
	}()

	// Wait for goroutines to finish
	wg.Wait()
	close(results)

	// Aggregate results
	for result := range results {
		if result {
			patternDetected = true
		}
	}

	return patternDetected
}

// detectPattern checks for patterns in a log file
func detectPattern(logFile string, pattern string, logger *logrus.Logger) bool {
	file, err := os.Open(logFile)
	if err != nil {
		logger.Errorf("Error opening log file %s: %v", logFile, err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("Error compiling regex pattern: %v", err)
		return false
	}

	patternDetected := false
	for scanner.Scan() {
		line := scanner.Text()
		if regex.MatchString(line) {
			logger.Warnf("Suspicious activity detected in %s: %s", logFile, line)
			patternDetected = true
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Errorf("Error reading log file %s: %v", logFile, err)
	}

	return patternDetected
}
