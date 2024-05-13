package cmd

import (
	"fmt"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
	"github.com/royiro10/cogo/server"
)

const LOCK_FILE = "./cogo.lock"
const (
	RUN_DAEMON     = "daemon"
	START_DAEMON   = "start"
	STOP_DAEMON    = "stop"
	RUN_COMMAND    = "run"
	KILL_COMMAND   = "kill"
	OUTPUT_COMMAND = "output"
	STATUS_COMMAND = "status"
	UNKOWN         = "unkown"
)

type CogoCLIDeps struct {
	Logger      *common.Logger
	LockService common.LockService
	CogoDaemon  server.Daemon
}

type CogoCLI struct {
	logger   *common.Logger
	commands map[string]models.CogoCLICommand
}

func CreateCLI(deps CogoCLIDeps) *CogoCLI {
	commands := map[string]models.CogoCLICommand{
		RUN_DAEMON:     makeHandleRunAsDaemon(deps.LockService, deps.Logger, deps.CogoDaemon),
		START_DAEMON:   makeHandleStartDaemon(deps.LockService, deps.Logger),
		STOP_DAEMON:    makeHandleStopDaemon(deps.LockService, deps.Logger),
		RUN_COMMAND:    makeRunCommand(deps.LockService, deps.Logger),
		KILL_COMMAND:   makeKillCommand(deps.LockService, deps.Logger),
		OUTPUT_COMMAND: makeOutputCommand(deps.LockService, deps.Logger),
		STATUS_COMMAND: makeStatusCommand(deps.LockService, deps.Logger),
	}

	service := &CogoCLI{
		commands: commands,
		logger:   deps.Logger,
	}

	return service
}

func (cli *CogoCLI) Handle(args []string, flags *models.CogoCLIFlags) {
	if len(args) < 1 {
		cli.UsageMsg()
		return
	}

	cmdInfo := models.CogoCLIInfo{
		Commad: args[0],
		Args:   args[1:],
		Flags:  flags,
	}

	if command := cli.commands[cmdInfo.Commad]; command != nil {
		err := command(cmdInfo)
		if err != nil {
			cli.logger.Fatal(err)
		}

		return
	}

	cli.logger.Info("unkown command has been run", "command", cmdInfo.Commad)

	fmt.Println("Unkown commnad: ", cmdInfo.Commad)
	cli.UsageMsg()
}

func (cli *CogoCLI) UsageMsg() {
	msg := `Usage: cogo [flags] [command] [arguments]

	Commands:
	  start                      Start the Cogo daemon.
	  stop                       Stop the Cogo daemon.
	  run "command"              Run a command in the background.
	  kill                       Kill a running command.
	  output                     Retrieve the output of a specific session.
	  status                     Get Cogo sessions' status
	
	Flags:
	  -s, --session <session-id> Specify a session to interact with.
	  -l, --logger               Enable logging.
	  -f, --follow               Stream the output of the running command.
	
Please refer to the README.md for more detailed information on each command and their usage.`

	fmt.Println(msg)
}
