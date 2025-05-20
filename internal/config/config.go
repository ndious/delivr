package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Discord    DiscordConfig `json:"discord" yaml:"discord"`
	Docker     DockerConfig  `json:"docker,omitempty" yaml:"docker,omitempty"`
	Logs       LogConfig     `json:"logs,omitempty" yaml:"logs,omitempty"`
	Commands   []Command     `json:"commands" yaml:"commands"`
	WorkingDir string        `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
}

// DiscordConfig holds Discord integration settings
type DiscordConfig struct {
	Token     string `json:"token" yaml:"token"`
	ChannelID string `json:"channelId" yaml:"channelId"`
}

// DockerConfig holds Docker-specific settings
type DockerConfig struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Directory string `json:"directory,omitempty" yaml:"directory,omitempty"`  // Directory to store log files
	MaxSize   int    `json:"maxSize,omitempty" yaml:"maxSize,omitempty"`    // Maximum size in MB before rotation
	MaxAge    int    `json:"maxAge,omitempty" yaml:"maxAge,omitempty"`     // Maximum age in days before deletion
	MaxBackups int   `json:"maxBackups,omitempty" yaml:"maxBackups,omitempty"` // Maximum number of backups to keep
	Compress  bool   `json:"compress,omitempty" yaml:"compress,omitempty"`   // Whether to compress rotated files
}

// Command represents a command to be executed
type Command struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description" yaml:"description"`
	Command     string   `json:"command" yaml:"command"`
	Args        []string `json:"args,omitempty" yaml:"args,omitempty"`
	Dir         string   `json:"dir,omitempty" yaml:"dir,omitempty"`
	EnvVars     []string `json:"envVars,omitempty" yaml:"envVars,omitempty"`
}

// Variables pour stocker le chemin du fichier de configuration chargé
var loadedConfigPath string

// DefaultConfigPath returns the default config file paths in order of preference
func DefaultConfigPath() string {
	// Try YAML first in current directory
	if _, err := os.Stat("config.yml"); err == nil {
		return "config.yml"
	}
	
	// Then try JSON in current directory
	if _, err := os.Stat("config.json"); err == nil {
		return "config.json"
	}
	
	// Then try in home directory
	home, err := os.UserHomeDir()
	if err == nil {
		// Try YAML in home directory
		homeYamlCfg := filepath.Join(home, ".delivr", "config.yml")
		if _, err := os.Stat(homeYamlCfg); err == nil {
			return homeYamlCfg
		}
		
		// Try JSON in home directory
		homeJsonCfg := filepath.Join(home, ".delivr", "config.json")
		if _, err := os.Stat(homeJsonCfg); err == nil {
			return homeJsonCfg
		}
	}
	
	// Default to current directory YAML
	return "config.yml"
}

// GetLoadedConfigPath returns the path of the loaded configuration file
func GetLoadedConfigPath() string {
	return loadedConfigPath
}

// isYAMLFile checks if a path has a YAML extension
func isYAMLFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yml" || ext == ".yaml"
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
	
	// Determine if it's a YAML file and use appropriate unmarshal
	if isYAMLFile(configPath) {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("error parsing YAML config: %w", err)
		}
	} else {
		// Assume JSON
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("error parsing JSON config: %w", err)
		}
	}
	
	// Store the loaded config path
	loadedConfigPath = configPath
	
	return &config, nil
}

// Save saves the configuration to file
func Save(config *Config, path string) error {
	var data []byte
	var err error
	
	// Determine format based on file extension
	if isYAMLFile(path) {
		data, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("error encoding YAML: %w", err)
		}
	} else {
		// Default to JSON
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("error encoding JSON: %w", err)
		}
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
