package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	l    *log.Logger
	once sync.Once
)

var serviceLogStyle = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 1, 0, 1).
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("#FBF6E9"))

func InitLogger(output *os.File, minLevel log.Level) {
	once.Do(func() {
		l = newLogger(output, minLevel)
	})
}

func newLogger(output *os.File, minLevel log.Level) *log.Logger {
	styles := log.DefaultStyles()

	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("63")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("86")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("192")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("203")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FATAL").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	logger := log.New(output)
	logger.SetStyles(styles)
	logger.SetLevel(minLevel)
	logger.SetTimeFormat(time.RFC3339)

	return logger
}

func Debug(service string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.Debug(fmt.Sprintf("%s %s", serviceLogStyle.Render(service), message))
}

func Info(service string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.Info(fmt.Sprintf("%s %s", serviceLogStyle.Render(service), message))
}

func Warn(service string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.Warn(fmt.Sprintf("%s %s", serviceLogStyle.Render(service), message))
}

func Error(service string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.Error(fmt.Sprintf("%s %s", serviceLogStyle.Render(service), message))
}

func Fatal(service string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.Fatal(fmt.Sprintf("%s %s", serviceLogStyle.Render(service), message))
}
