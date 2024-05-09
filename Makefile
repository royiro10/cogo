PROJECT_NAME=cogo

BUILD_FLAGS=
MAIN_SRC_FILE=./main.go
OUT_EXEC_FILE=./$(PROJECT_NAME).exe

all: clear build run


clear:
	-$(OUT_EXEC_FILE) stop

# to allow for gracefull shutdown
	sleep 1 
	
	-rm $(OUT_EXEC_FILE) ./$(PROJECT_NAME).lock
	-find . -wholename "$($(PROJECT_NAME))*.log" -delete


build: 
	go build -o $(OUT_EXEC_FILE) $(MAIN_SRC_FILE)


run:
	$(OUT_EXEC_FILE) start



###################
#   dev scripts   #
###################

#get the pid from the lock file and match it to the logs files
DAEMON_PID=$(shell cat ./cogo.lock | grep -Po '"Pid":\K\d+')
DAEMON_LOG_FILE=$(shell find ./logs -name "$(PROJECT_NAME)_$(DAEMON_PID).log")

END_COLOR=\033[0m
COLOR_RED=\033[0;31m
COLOR_GREEN=\033[0;32m
COLOR_YELLOW=\033[0;32m

daemon_logs:
	@tail -f ${DAEMON_LOG_FILE} | while read line; \
	do \
		if echo "$${line}" | grep -q "level=ERROR"; then \
			echo -e "$(COLOR_RED)level=$${line#*level=}$(END_COLOR)"; \
		elif echo "$${line}" | grep -q "level=DEBUG"; then \
			echo -e "$(COLOR_GREEN)level=$${line#*level=}$(END_COLOR)"; \
		elif echo "$${line}" | grep -q "level=WARN"; then \
			echo -e "$(COLOR_YELLOW)level=$${line#*level=}$(END_COLOR)"; \
		else \
			echo "level=$${line#*level=}"; \
		fi; \
	done \




