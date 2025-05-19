package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ndious/delivr/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// CommandLogger is responsible for logging command output to files
type CommandLogger struct {
	config  config.LogConfig
	baseDir string
	loggers map[string]*lumberjack.Logger
}

// NewCommandLogger creates a new command logger
func NewCommandLogger(cfg config.LogConfig) (*CommandLogger, error) {
	// Set default values if not specified
	if cfg.Directory == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		cfg.Directory = filepath.Join(homeDir, ".delivr", "logs")
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = 10 // 10 MB
	}

	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30 // 30 days
	}

	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 5
	}

	// Ensure log directory exists
	if err := os.MkdirAll(cfg.Directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	return &CommandLogger{
		config:  cfg,
		baseDir: cfg.Directory,
		loggers: make(map[string]*lumberjack.Logger),
	}, nil
}

// GetLogWriter returns a writer for the specified command
func (l *CommandLogger) GetLogWriter(commandName string) io.Writer {
	// Sanitize command name for use in filenames
	safeCommandName := sanitizeFilename(commandName)

	// Check if logger already exists
	if logger, ok := l.loggers[safeCommandName]; ok {
		return logger
	}

	// Create new logger
	today := time.Now().Format("2006-01-02")
	logPath := filepath.Join(l.baseDir, fmt.Sprintf("%s-%s.log", safeCommandName, today))

	logger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    l.config.MaxSize,
		MaxBackups: l.config.MaxBackups,
		MaxAge:     l.config.MaxAge,
		Compress:   l.config.Compress,
	}

	l.loggers[safeCommandName] = logger
	return logger
}

// GetLogPath returns the log file path for a command
func (l *CommandLogger) GetLogPath(commandName string) string {
	safeCommandName := sanitizeFilename(commandName)
	today := time.Now().Format("2006-01-02")
	return filepath.Join(l.baseDir, fmt.Sprintf("%s-%s.log", safeCommandName, today))
}

// Close closes all open loggers
func (l *CommandLogger) Close() {
	for _, logger := range l.loggers {
		_ = logger.Close()
	}
}

// sanitizeFilename removes characters that are problematic in filenames
func sanitizeFilename(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, ":", "-")
	name = strings.ReplaceAll(name, "*", "-")
	name = strings.ReplaceAll(name, "?", "-")
	name = strings.ReplaceAll(name, "\"", "-")
	name = strings.ReplaceAll(name, "<", "-")
	name = strings.ReplaceAll(name, ">", "-")
	name = strings.ReplaceAll(name, "|", "-")
	return name
}
