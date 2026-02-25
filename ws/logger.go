package ws

import "github.com/atareversei/http-server/internal/cli"

// DefaultLogger provides a basic implementation of the Logger interface using the cli package.
type DefaultLogger struct{}

func (dl *DefaultLogger) Success(msg string) {
	cli.Success(msg)
}
func (dl *DefaultLogger) Info(msg string) {
	cli.Info(msg)
}
func (dl *DefaultLogger) Warning(msg string) {
	cli.Warning(msg)
}
func (dl *DefaultLogger) Error(msg string, err error) {
	cli.Error(msg, err)
}

// Logger defines the interface for logging messages at various levels.
type Logger interface {
	Success(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string, err error)
}

// NoOpLogger implements Logger but discards all log output.
// Useful for disabling logs without changing logic.
type NoOpLogger struct{}

func (n *NoOpLogger) Success(msg string)          {}
func (n *NoOpLogger) Info(msg string)             {}
func (n *NoOpLogger) Warning(msg string)          {}
func (n *NoOpLogger) Error(msg string, err error) {}

// DisableLogging disables all logging by replacing the logger with a no-op logger.
func (c *Connection) DisableLogging() {
	c.previousLogger = c.Logger
	c.Logger = &NoOpLogger{}
	c.loggingEnabled = false
}

// EnableLogging restores the previously used logger and re-enables logging.
func (c *Connection) EnableLogging() {
	c.Logger = c.previousLogger
	c.loggingEnabled = true
}
