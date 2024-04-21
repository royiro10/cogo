package util

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

var DefaultLogger = CreateLogger(fmt.Sprintf("./logs/cogo_%d.log", os.Getpid()))

type Logger struct {
	*slog.Logger
}

func CreateLogger(logFile string) *Logger {
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logger := &Logger{
			slog.Default(),
		}

		LogErrorFatal(logger, "could not create log file", err)
		return nil
	}
	// TODO: should I really close the logging file?
	// defer f.Close()

	writer := io.MultiWriter(os.Stdout, f)
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	return &Logger{
		slog.New(handler),
	}
}

func LogErrorFatal(l *Logger, msg string, err error) {
	panic(fmt.Errorf("%s %w", msg, err).Error())
}
