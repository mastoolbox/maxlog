
# maxlog

## Overview

`maxlog` is a command-line tool that retrieves and displays logs from the Maximo Application Suite (MAS). It supports multiple operating modes and provides straightforward commands for inspection and troubleshooting.
## Installation

```bash
go build
```

This will generate the `maxlog` executable.

## Usage

```bash
maxlog [action] [options]
```

### Actions

- `logs` – Show container or pod logs
- `inspect` – Inspect pods or containers
- `version` – Display the current version
- `help` – Show help information

### Example

```bash
maxlog logs --tag=mytag --tail=100
```

## Environment Variables

- `MAXLOG_MODE`  
  Sets the operation mode (`k8s` for Kubernetes, `pod` for Podman)
- `MAXLOG_TAIL`  
  Number of log lines to display (default: 40)
- `MAXLOG_K8S_NAMESPACE`  
  Namespace for Kubernetes logs
- `MAXLOG_K8S_APPTYPE` - optional  
  This is the pod selector for Kubernetes logs. The default value is `all`, `ui`, `cron`, `mea`, `rpt`, `jms`. This is only required in k8s mode.
- `MAXLOG_CONTAINER`  
  Container name in Podman mode
- `MAXLOG_USE_NERDFONT` - optional  
  With the values `1` or `true`, symbols can be activated via a NerdFont (see [Nerd Fonts](https://www.nerdfonts.com/)).  This option is disabled by default.
- `MAXLOG_FOCUS` - optional  
  It hides all lines that do not contain the word. It is not case-sensitive.

### Configuration file example

```bash
#!/usr/bin/env bash
export MAXLOG_MODE="k8s"
export MAXLOG_TAIL="100"
export MAXLOG_K8S_NAMESPACE="mas-demo-manage"
export MAXLOG_K8S_APPTYPE="all"
export MAXLOG_USE_NERDFONT="true"
```

## Use
After setting the mode and selector via environment variables, `maxlog` can be used without restrictions. See the section on environment variables for more information. Without Command, the logs command is used:
```bash
maxlog
```
Only debug information should be displayed and one tag highlighted:
```bash
maxlog logs focus=debug tag=ZZTEST
```
The flags tag and focus make the command logs unnecessary:
```bash
maxlog focus=debug tag=ZZTEST
```

## Building
In the Go programming language, the following command is generally executed within the cloned directory:
```bash
go build
```
This will provide an example of cross-compiling a release:
```bash
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w"
```

## License

See the `LICENSE` file in the project.
