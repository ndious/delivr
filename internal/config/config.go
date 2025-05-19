package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the main configuration structure
type Config struct {
	Discord    DiscordConfig `json:"discord"`
	Docker     DockerConfig  `json:"docker,omitempty"`
	Logs       LogConfig     `json:"logs,omitempty"`
	Commands   []Command     `json:"commands"`
	WorkingDir string        `json:"workingDir,omitempty"`
}

// DiscordConfig holds Discord integration settings
type DiscordConfig struct {
	Token     string `json:"token"`
	ChannelID string `json:"channelId"`
}

// DockerConfig holds Docker-specific settings
type DockerConfig struct {
	Host string `json:"host,omitempty"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Directory string `json:"directory,omitempty"`  // Directory to store log files
	MaxSize   int    `json:"maxSize,omitempty"`    // Maximum size in MB before rotation
	MaxAge    int    `json:"maxAge,omitempty"`     // Maximum age in days before deletion
	MaxBackups int   `json:"maxBackups,omitempty"` // Maximum number of backups to keep
	Compress  bool   `json:"compress,omitempty"`   // Whether to compress rotated files
}

// Command represents a command to be executed
type Command struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Command     string   `json:"command"`
	Args        []string `json:"args,omitempty"`
	Dir         string   `json:"dir,omitempty"`
	EnvVars     []string `json:"envVars,omitempty"`
}

// Variables pour stocker le chemin du fichier de configuration chargé
var loadedConfigPath string

// DefaultConfigPath returns the default config file path
func DefaultConfigPath() string {
	// First try current directory
	if _, err := os.Stat("config.json"); err == nil {
		return "config.json"
	}
	
	// Then try home directory
	home, err := os.UserHomeDir()
	if err == nil {
		homeCfg := filepath.Join(home, ".delivr", "config.json")
		if _, err := os.Stat(homeCfg); err == nil {
			return homeCfg
		}
	}
	
	// Default to current directory anyway
	return "config.json"
}

// GetLoadedConfigPath returns the path of the loaded configuration file
func GetLoadedConfigPath() string {
	return loadedConfigPath
}

// Load loads the configuration from file
func Load(customPath string) (*Config, error) {
	configPath := DefaultConfigPath()
	
	// Check if config path is provided as a parameter
	if customPath != "" {
		configPath = customPath
	} else if envPath := os.Getenv("DELIVR_CONFIG"); envPath != "" {
		// Check if config path is overridden by environment
		configPath = envPath
	}
	
	// Vérifier que le fichier existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	
	// Store the loaded config path
	loadedConfigPath = configPath
	
	return &config, nil
}

// Save saves the configuration to file
func Save(config *Config, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig(path string) error {
	// Create a default configuration
	defaultConfig := &Config{
		WorkingDir: "",
		Docker: DockerConfig{
			Host: "unix:///var/run/docker.sock",
		},
		Logs: LogConfig{
			Directory: "./logs",
			MaxSize:   10,
			MaxAge:    30,
			MaxBackups: 5,
			Compress:  true,
		},
		Discord: DiscordConfig{
			Token:     "YOUR_DISCORD_BOT_TOKEN_HERE",
			ChannelID: "YOUR_DISCORD_CHANNEL_ID_OR_WEBHOOK_URL_HERE",
		},
		Commands: []Command{
			{
				Name:        "Show Docker Status",
				Description: "Lists all running Docker containers",
				Command:     "docker",
				Args:        []string{"ps", "-a"},
			},
			{
				Name:        "Git Status",
				Description: "Shows the working tree status",
				Command:     "git",
				Args:        []string{"status"},
			},
		},
	}
	
	// Save the configuration to the specified path
	return Save(defaultConfig, path)
}
