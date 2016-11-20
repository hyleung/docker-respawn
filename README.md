# Docker-respawn

Small PoC program that watches for the HEALTCHECK status of a container and "respawns" the container
if it is deemed to be unhealthy.

## Usage

```
NAME:
   docker-respawn - Restart Docker containers that fail health-check

USAGE:
   docker-respawn <image name>

VERSION:
   0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        Enable debug level logging
   --help, -h     show help
   --version, -v  print the version
```
