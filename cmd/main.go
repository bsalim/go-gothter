package main

import (
	"bufio"
	"go-gothter/internal/config"
	"go-gothter/internal/notifier"
	"go-gothter/internal/scanner"
	"go-gothter/internal/utils"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/default_config.yaml")
	if err != nil {
		logrus.Fatalf("Error loading config: %v", err)
	}

	// Convert block duration to time.Duration
	blockDuration := time.Duration(cfg.BlockDuration.Hours) * time.Hour

	// Create or open the log file
	logFile, err := os.OpenFile(cfg.LogFiles.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
	logger.SetOutput(logFile) // Set log output to the file

	// Initialize notifier
	logger.Info("Starting Gothter...")
	not := notifier.NewNotifier(cfg, logger)

	// Initialize log scanner
	scanner := scanner.NewLogScanner(cfg, logger)

	// Start monitoring log files
	go scanner.StartMonitoring()

	// Monitor log files
	for {
		patternDetected, ip := monitorLogs(cfg, logger, blockDuration)
		if patternDetected && ip != "" {
			not.SendNotification(ip)
		}

		// Sleep for a while before checking again
		time.Sleep(5 * time.Minute)
	}
}

// monitorLogs monitors log files concurrently
func monitorLogs(cfg *config.Config, logger *logrus.Logger, blockDuration time.Duration) (bool, string) {
	var wg sync.WaitGroup
	patternDetected := false
	var detectedIP string
	results := make(chan struct {
		detected bool
		ip       string
	}, 2) // Channel to collect results from goroutines

	// Monitor auth.log
	wg.Add(1)
	go func() {
		defer wg.Done()
		ip := ""
		if detectPattern(cfg.LogFiles.AuthLog, cfg.Patterns.AuthFail, blockDuration, logger, &ip) {
			results <- struct {
				detected bool
				ip       string
			}{true, ip}
		}
	}()

	// Monitor nginx.log
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	ip := ""
	// 	if detectPattern(cfg.LogFiles.NginxLog, cfg.Patterns.Nginx404, blockDuration, logger, &ip) {
	// 		results <- struct {
	// 			detected bool
	// 			ip       string
	// 		}{true, ip}
	// 	}
	// }()

	// Wait for goroutines to finish
	wg.Wait()
	close(results)

	// Aggregate results
	for result := range results {
		if result.detected {
			patternDetected = true
			detectedIP = result.ip
		}
	}

	return patternDetected, detectedIP
}

// detectPattern checks for patterns in a log file
func detectPattern(logFile string, pattern string, blockDuration time.Duration, logger *logrus.Logger, detectedIP *string) bool {
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
			ip := utils.ExtractIP(line)
			if ip != "" {
				utils.BlockIP(ip, blockDuration, logger)
				*detectedIP = ip
				patternDetected = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Errorf("Error reading log file %s: %v", logFile, err)
	}

	return patternDetected
}
