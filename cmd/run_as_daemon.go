package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/server"
)

func makeHandleRunAsDaemon(lockService common.LockService, logger *common.Logger, daemon server.Daemon) models.CogoCLICommand {
	handleDaemon := func(cmdInfo models.CogoCLIInfo) error {
		if lockService.IsAcquired(LOCK_FILE) {
			logger.Info("Daemon is already running.")
		}

		release, err := lockService.Acquire(LOCK_FILE)
		defer release()
		if err != nil {
			logger.Error("can not acquire lock", "err", err)
			return fmt.Errorf("can not acquire lock: %w", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var stopChan = make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			sig := <-stopChan
			logger.Info("Received interrupt signal, stopping...", "signal", sig.String())
			cancel()
			release()
			os.Exit(1)
		}()

		daemon.Start(ctx)
		return nil
	}

	return handleDaemon
}
