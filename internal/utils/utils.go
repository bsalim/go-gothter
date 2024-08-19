package utils

import (
	"os/exec"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

// BlockIP blocks the IP address using iptables
func BlockIP(ip string, duration time.Duration, logger *logrus.Logger) {
	cmd := exec.Command("sudo", "iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
	err := cmd.Run()
	if err != nil {
		logger.Errorf("Error blocking IP %s: %v", ip, err)
		return
	}
	logger.Infof("Blocked IP %s", ip)

	// Schedule unblock using the duration from config
	go ScheduleUnblockIP(ip, duration, logger)
}

// ScheduleUnblockIP schedules the IP unblock after a specified duration
func ScheduleUnblockIP(ip string, duration time.Duration, logger *logrus.Logger) {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	// Wait for the timer to expire
	<-timer.C
	UnblockIP(ip, logger)
}

// UnblockIP unblocks the IP address using iptables
func UnblockIP(ip string, logger *logrus.Logger) {
	cmd := exec.Command("sudo", "iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
	err := cmd.Run()
	if err != nil {
		logger.Errorf("Error unblocking IP %s: %v", ip, err)
		return
	}
	logger.Infof("Unblocked IP %s", ip)
}

// ExtractIP extracts the IP address from a log line using a regex
func ExtractIP(line string) string {
	// Regex pattern to match IP addresses
	ipRegex := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	match := ipRegex.FindStringSubmatch(line)
	if len(match) > 0 {
		return match[0]
	}
	return ""
}
