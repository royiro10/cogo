package commands

import (
	"bufio"
	"io"
	"time"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
)

func readPipe(pipe io.ReadCloser, output chan<- models.StdLine, cwd string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		if msg := scanner.Text(); msg != "" {
			output <- models.StdLine{
				Cwd:  cwd,
				Time: time.Now(),
				Data: msg,
			}
		}
	}

	pipe.Close()
}

func makePipeLogger(sc *StdContainer, logger *common.Logger) *StdListener {
	var pipeLoggerListener StdListener = func(line *models.StdLine) {
		logger.Info(line.Data, "pipe", sc.Name)
	}

	return &pipeLoggerListener
}
