package common

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// var DefaultLogger = CreateLogger(fmt.Sprintf("./logs/cogo_%d.log", os.Getpid()))
var EmptyLogger = &Logger{slog.New(slog.NewJSONHandler(io.Discard, nil)), nil}

type Logger struct {
	*slog.Logger
	LogFile *os.File
}

type LoggerOptions struct {
	LogPath string
	LogFile string
	Level   *slog.Level
}

func CreateLogger(options *LoggerOptions) *Logger {
level := slog.LevelInfo
	if len(levels) > 0 {
		level = levels[0]
	}
if _, err := os.Stat(options.LogPath); os.IsNotExist(err) {
		err := os.MkdirAll(options.LogPath, 0644)
		if err != nil {
			panic(fmt.Sprintf("Could not make directory %q", options.LogPath))
		}
	}
	f, err := os.OpenFile(
		filepath.Join(options.LogPath, options.LogFile),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		logger := &Logger{
			slog.Default(),
			nil,
		}

		logger.Fatal(WrapedError{"could not create log file", err})
		return nil
	}
	// TODO: should I really close the logging file?
	// defer f.Close()

	writer := io.MultiWriter(os.Stdout, f)
	handler := slog.NewTextHandler(writer, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})
	return &Logger{
		slog.New(handler),
		f,
	}
}

type WrapedError struct {
	Msg string
	Err error
}

func (l *Logger) Fatal(input interface{}) {
	var err error
	switch v := input.(type) {
	case string:
		err = fmt.Errorf(v)
	case error:
		err = v
	case WrapedError:
		err = fmt.Errorf("%s, err: %s", v.Msg, v.Err.Error())
	default:
		l.Warn("Received an unsupported type: %T\n", v)
	}

	l.Error(err.Error())
	panic(err)
}
