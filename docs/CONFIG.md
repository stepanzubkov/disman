# Disman Configuration

Configuration are provided through CLI options and config file.

You can see CLI options by running `disman --help`

## Configuration file

Configuration file basically stored in `/etc/disman.conf`. It handles as `.desktop` file (desktop entry)

You can find sample configuration file in this repo (`res/conf`)

Configuration options:

`DISPLAY` (string) - X11 display name. Defaults to `:0`. Should be in format `[host]:(display).[screen]`

`VT` (string) - Virtual terminal name (TTY). Should be in format `vtX`, where `X` is VT's number. Defaults to `vt7` (7th virtual terminal)

`PRE_COMMAND` (command) - Command should be run before the asking user for username/password and before printing title.

`DISPLAY_TITLE` (bool) - If true, enables printing title (`>>> Disman Display Manager <<<`) at the start of login.
