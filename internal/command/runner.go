package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ndious/delivr/internal/config"
)

// Discord interface defines the methods required for discord integration
type Discord interface {
	SendMessage(content string) error
}

// Logger interface defines the methods required for logging
type Logger interface {
	GetLogWriter(commandName string) io.Writer
	GetLogPath(commandName string) string
}

// Runner executes commands
type Runner struct {
	discord    Discord
	logger     Logger
	workingDir string
	dockerHost string
}

// NewRunner creates a new command runner
func NewRunner(discord Discord, logger Logger, workingDir string, dockerHost string) *Runner {
	return &Runner{
		discord:    discord,
		logger:     logger,
		workingDir: workingDir,
		dockerHost: dockerHost,
	}
}

// Execute runs a command and sends its output to Discord
func (r *Runner) Execute(cmd config.Command) error {
	startTime := time.Now()

	// Prepare notification message
	startMsg := fmt.Sprintf("🏃 Running command: **%s**\n> %s", cmd.Name, cmd.Description)
	if err := r.discord.SendMessage(startMsg); err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	// Prepare command
	command := exec.Command(cmd.Command, cmd.Args...)

	// Set Docker host if specified
	if r.dockerHost != "" && cmd.Command == "docker" {
		env := os.Environ()
		env = append(env, "DOCKER_HOST="+r.dockerHost)
		command.Env = env
	}

	// Set working directory based on priority:
	// 1. Command-specific directory if specified
	// 2. Global working directory if specified
	// 3. Current directory otherwise
	if cmd.Dir != "" {
		command.Dir = cmd.Dir
	} else if r.workingDir != "" {
		command.Dir = r.workingDir
	}

	// Set environment variables if specified
	if len(cmd.EnvVars) > 0 {
		command.Env = append(os.Environ(), cmd.EnvVars...)
	}

	// Get log writer for this command
	logWriter := r.logger.GetLogWriter(cmd.Name)

	// Create multi-writers to capture output in memory and log to file
	var stdout, stderr bytes.Buffer
	multiStdout := io.MultiWriter(&stdout, logWriter)
	multiStderr := io.MultiWriter(&stderr, logWriter)

	// Write command metadata to log file
	fmt.Fprintf(logWriter, "\n\n==================================================\n")
	fmt.Fprintf(logWriter, "Command: %s\n", cmd.Name)
	fmt.Fprintf(logWriter, "Description: %s\n", cmd.Description)
	fmt.Fprintf(logWriter, "Executed at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(logWriter, "Working Directory: %s\n", command.Dir)
	fmt.Fprintf(logWriter, "Full Command: %s %s\n", cmd.Command, strings.Join(cmd.Args, " "))
	fmt.Fprintf(logWriter, "==================================================\n\n")

	// Set output writers
	command.Stdout = multiStdout
	command.Stderr = multiStderr

	// Execute the command
	err := command.Run()

	// Log completion status
	if err != nil {
		fmt.Fprintf(logWriter, "\n\n==================================================\n")
		fmt.Fprintf(logWriter, "Command failed with error: %v\n", err)
		fmt.Fprintf(logWriter, "==================================================\n\n")
	} else {
		fmt.Fprintf(logWriter, "\n\n==================================================\n")
		fmt.Fprintf(logWriter, "Command completed successfully\n")
		fmt.Fprintf(logWriter, "==================================================\n\n")
	}

	// Calculate execution time
	duration := time.Since(startTime)
	durationStr := fmt.Sprintf("%.2f seconds", duration.Seconds())

	// Prepare output for Discord
	var resultMsg strings.Builder
	if err != nil {
		resultMsg.WriteString(fmt.Sprintf("❌ Command **%s** failed (took %s)\n", cmd.Name, durationStr))
		if stderr.Len() > 0 {
			errText := stderr.String()
			// Truncate if too long
			if len(errText) > 1500 {
				errText = errText[:1500] + "... (truncated)"
			}
			resultMsg.WriteString(fmt.Sprintf("```\n%s\n```", errText))
		} else {
			resultMsg.WriteString(fmt.Sprintf("Error: %v", err))
		}
	} else {
		resultMsg.WriteString(fmt.Sprintf("✅ Command **%s** completed successfully (took %s)\n", cmd.Name, durationStr))
		if stdout.Len() > 0 {
			outText := stdout.String()
			// Truncate if too long
			if len(outText) > 1500 {
				outText = outText[:1500] + "... (truncated)"
			}
			resultMsg.WriteString(fmt.Sprintf("```\n%s\n```", outText))
		}
	}

	// Add log file info to result
	logPath := r.logger.GetLogPath(cmd.Name)
	resultMsg.WriteString(fmt.Sprintf("\n📄 Log file: `%s`", logPath))

	// Send result to Discord
	if err := r.discord.SendMessage(resultMsg.String()); err != nil {
		return fmt.Errorf("failed to send result message: %w", err)
	}

	return err
}

// ExecuteAll runs all commands in sequence
func (r *Runner) ExecuteAll(commands []config.Command) error {
	for _, cmd := range commands {
		err := r.Execute(cmd)
		if err != nil {
			return fmt.Errorf("command '%s' failed: %w", cmd.Name, err)
		}
	}
	return nil
}
