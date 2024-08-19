package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	LogFiles struct {
		AuthLog  string `yaml:"auth_log"`
		NginxLog string `yaml:"nginx_log"`
	} `yaml:"log_files"`
	Patterns struct {
		AuthFail string `yaml:"auth_fail"`
		Nginx404 string `yaml:"nginx_404"`
	} `yaml:"patterns"`
	Email struct {
		Enabled      bool   `yaml:"enabled"`
		SMTPServer   string `yaml:"smtp_server"`
		SMTPPort     int    `yaml:"smtp_port"`W
		SMTPUser     string `yaml:"smtp_user"`
		SMTPPassword string `yaml:"smtp_password"`
		Recipient    string `yaml:"recipient"`
		Subject      string `yaml:"subject"`
	} `yaml:"email"`
	// Slack notification coming soon
	// Slack struct {
	// 	Enabled    bool   `yaml:"enabled"`
	// 	WebhookURL string `yaml:"webhook_url"`
	// 	Channel    string `yaml:"channel"`
	// 	Username   string `yaml:"username"`
	// } `yaml:"slack"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	return &cfg, nil
}
