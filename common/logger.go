package common

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

// var DefaultLogger = CreateLogger(fmt.Sprintf("./logs/cogo_%d.log", os.Getpid()))
var EmptyLogger = &Logger{slog.New(slog.NewJSONHandler(io.Discard, nil)), nil}

type Logger struct {
	*slog.Logger
	LogFile *os.File
}

func CreateLogger(logFile string) *Logger {
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
		Level:     slog.LevelDebug,
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
