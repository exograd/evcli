% evcli

# Introduction
The evcli program interacts with the [Eventline
platform](https://eventline.net). It can be used to retrieve information about
various components, manage projects, control pipelines…

# Conventions
## Status
All commands exit with a status code equal to zero if they succeed. On
failure, a non-zero status code is used. The status code 1 is the only error
status code used for the time being.

## Output
Commands write the output of their main operation to stdout. All other content
is written to stderr; this includes all kinds of messages and data table
headers.

As a result, the output of evcli can always be processed by various tools:
you can simply pipe evcli commands into other programs without having to
filter out any content.

## Information
Commands usually do not display any status or information message when they
succeed, in accordance with the UNIX Rule of Silence.

The `--verbose` global option can be used to enable informational messages.

# Configuration
Configuration is stored in a regular file at `$HOME/.evcli/config.json`. If
the file does not exist, evcli will create it and fill it with the default
configuration.

Note that the configuration file will store your API key. If you write the
configuration file yourself, make sure to set permissions to `0600`. Evcli
checks for permissions at startup and aborts if the configuration file is
readable for anyone but the file owner.

## Settings
Settings are stored hierarchically. For example, for `api.endpoint`, the
`endpoint` member is stored in an `api` top-level JSON object.

The following settings are available:

- `api.endpoint`: the base URI of the Eventline API server. This setting
should not be modified.
- `api.key`: the API key used by evcli to authenticate on the API server.
- `interface.color`: whether to activate color in the output of evcli or not.
  This setting is overridden by the `--color` option.

The `show-config`, `get-config` and `set-config` commands can be used to
interact with the current configuration.

# Options
## Global options
The following options can be used with all commands:

- `--debug <level>`: enable debug messages of level `<level>` or higher. The
  level is a positive integer; higher log levels are used for more detailed
  debug messages.
- `-h`, `--help`: print information about a command, its options and its
  arguments.
- `--no-color`: do not use any color in the output.
- `-v`, `--verbose`: enable status and information messages.
- `-y`, `--yes`: skip all confirmation prompts, automatically using "yes" as
  response.
  
## Project selection options
Most commands are executed in the context of a specific [Eventline
project](https://doc.eventline.net/organization/projects). By default, evcli
will look for a project file, `eventline-project.json` in the current
directory, and use it to identify the current project. Alternatively, you can
use one of the project selection options to choose the current project:

- `-p <name>`, `--project-name <name>`: select the current project by name.
- `--project-id <id>`: select the current project by identifier.
- `--project-path <path>`: select the current project using the path of a
  project directory; this directory is expected to contain an
  `eventline-project.json` project file.
  
Project selection options can be used with all commands. They will be ignored
for commands which do not depend on a current project.

# Commands
**TODO**
