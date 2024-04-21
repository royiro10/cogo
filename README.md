Cogo
====

Cogo is a powerful command-line tool that simplifies the management of terminal commands by running them as a daemon process in the background. It allows you to execute commands in a "shoot and forget" style, enabling you to send commands without maintaining multiple terminal windows open. This tool is especially useful for managing long-running tasks and reviewing their output at your convenience.

Features
--------

* **Background Execution:** Run commands in the background without keeping terminal sessions open.
* **Command Management:** Easily start, stop, and monitor the status of commands.
* **Interactive Retrieval:** Retrieve and interact with the output of your commands whenever needed.
* **Centralized Control:** Manage all your terminal commands from a single access point.

Getting Started
---------------

### Installation

Clone the repository and build the executable:

```bash
git clone https://github.com/royiro10/cogo.git 
cd cogo
go build -o cogo
```

try building it with make, there will be some utilities for those needind it there.

### Running Cogo

To start the Cogo daemon:

```bash
./cogo start
```

Sending a command to the daemon:

```bash
./cogo run "your-command-here"
```

To view running commands:

```bash
./cogo status
```

To retrieve the output of a specific command:

```bash
./cogo output <command-id>
```

To stop the daemon:

```bash
./cogo stop
```

Configuration
-------------

You can configure Cogo by modifying the `config.yaml` file located in the root directory. Available configurations include log file paths, maximum number of concurrent commands, and more.

Contributing
------------

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feat/amazing-feature`)
5. Open a Pull Request

License
-------

Distributed under the MIT License. See `LICENSE` for more information.
