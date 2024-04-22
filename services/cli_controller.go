package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/util"
)

const LOCK_FILE = "./cogo.lock"
const (
	RUN_DAEMON   = "daemon"
	START_DAEMON = "start"
	STOP_DAEMON  = "stop"
	RUN_COMMAND  = "run"
	UNKOWN       = "unkown"
)

type CogoCLIDeps struct {
	Logger      *util.Logger
	LockService LockService
	CogoDaemon  Daemon
}

type CogoCLI struct {
	logger   *util.Logger
	commands map[string]models.CogoCLICommand
}

func CreateCLI(deps CogoCLIDeps) *CogoCLI {
	commands := map[string]models.CogoCLICommand{
		RUN_DAEMON:   makeHandleRunAsDaemon(deps.LockService, deps.Logger, deps.CogoDaemon),
		START_DAEMON: makeHandleStartDaemon(deps.LockService, deps.Logger),
		STOP_DAEMON:  makeHandleStopDaemon(deps.LockService, deps.Logger),
		RUN_COMMAND:  makeRunCommand(deps.LockService, deps.Logger),
	}

	service := &CogoCLI{
		commands: commands,
	}

	return service
}

func (cli *CogoCLI) Handle() {
	if len(os.Args) < 2 {
		cli.UsageMsg()
		return
	}

	cmdInfo := models.CogoCLIInfo{
		Commad: os.Args[1],
		Args:   os.Args[1:],
	}

	if command := cli.commands[cmdInfo.Commad]; command != nil {
		err := command(cmdInfo)
		if err != nil {
			util.LogErrorFatal(cli.logger, "", err)
		}

		return
	}

	cli.logger.Info("unkown command has been run", "command", cmdInfo.Commad)

	fmt.Println("Unkown commnad: ", cmdInfo.Commad)
	cli.UsageMsg()
}

func (cli *CogoCLI) UsageMsg() {
	avilableCommands := make([]string, len(cli.commands))

	i := 0
	for k := range cli.commands {
		avilableCommands[i] = k
		i++
	}

	fmt.Printf("Usage: cogo <%s>", strings.Join(avilableCommands, "|"))
}

func makeHandleRunAsDaemon(lockService LockService, logger *util.Logger, daemon Daemon) models.CogoCLICommand {
	handleDaemon := func(cmdInfo models.CogoCLIInfo) error {
		if lockService.IsAquired(LOCK_FILE) {
			logger.Info("Daemon is already running.")
		}

		release, err := lockService.Acquire(LOCK_FILE)
		defer release()
		if err != nil {
			logger.Error("can not aquire lock", "err", err)
			return fmt.Errorf("can not aquire lock: %w", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var stopChan = make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			<-stopChan
			logger.Info("Received interrupt signal, stopping...")
			cancel()
		}()

		daemon.Start(ctx)
		return nil
	}

	return handleDaemon
}

func makeHandleStartDaemon(lockService LockService, logger *util.Logger) models.CogoCLICommand {
	return func(cci models.CogoCLIInfo) error {
		_, err := startDaemon(lockService, logger)
		return err
	}
}

func startDaemon(lockService LockService, logger *util.Logger) (bool, error) {
	if lockService.IsAquired(LOCK_FILE) {
		logger.Info("Daemon is already running.")
		return false, nil
	}

	logger.Info("Starting daemon...", "arg0", os.Args[0])

	cmd := exec.Command(os.Args[0], RUN_DAEMON)
	err := cmd.Start()
	if err != nil {
		return false, fmt.Errorf("cmd.Start failed: %w ", err)
	}

	err = cmd.Process.Release()
	if err != nil {
		return false, fmt.Errorf("cmd.Process.Release failed: %w ", err)
	}

	return true, nil
}

func makeHandleStopDaemon(lockService LockService, logger *util.Logger) models.CogoCLICommand {
	handleStopDaemon := func(cmdInfo models.CogoCLIInfo) error {
		lockCommit, err := lockService.GetLockCommit(LOCK_FILE)
		if err != nil {
			return fmt.Errorf("can not get commited lock: %w", err)
		}

		if err = util.SendInterrupt(lockCommit.Pid); err != nil {
			logger.Error("could not send interrupt to process", "err", err, "pid", lockCommit.Pid)
			logger.Warn("hard kill", "pid", lockCommit.Pid)
			if err := util.KillCmd(lockCommit.Pid).Start(); err != nil {
				return fmt.Errorf("hard kill failed: %w", err)
			}
		}

		return lockService.Release(LOCK_FILE)
	}

	return handleStopDaemon
}

func makeRunCommand(lockService LockService, logger *util.Logger) models.CogoCLICommand {
	return func(cmdInfo models.CogoCLIInfo) error {
		if !lockService.IsAquired(LOCK_FILE) {
			// TODO: automaticcly start daemon
			return fmt.Errorf("cogo must be start before running commands")
		}

		client := CreateCogoClient(logger)
		defer client.releaseFunc()

		commandRequest := &CommandParameters{
			SessionId: DefaultSessionKey,
			Command:   strings.Join(cmdInfo.Args[1:], " "),
		}

		client.Run(commandRequest)
		return nil
	}
}
