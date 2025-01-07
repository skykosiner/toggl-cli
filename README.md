# toggl-cli
* A small CLI tool to control [Toggl](https://toggl.com) using the rest API
## Why?
* I really like to track my time using Toggl but I find the website and mobile app clunky and slow
* For a while I used [timery for toggl](https://timeryapp.com/), but I use Android and Linux now not iOS and Mac
    * But I really liked the way the app gave you saved timers so I took tho concept and made my own CLI tool to do the same
## Configuration
* The basic config is quite simple
`~/.config/toggl/config.json`
```json
{
   "api_key": "your api key",
   "workspace_id": your workspace id,
   "saved_timers": []
 }
```
* Ignore the saved timers array for now as that gets populated when you tell your program to save a new timer
## How to use
Running `toggl-cli` with no args will provide a basic help menu:
```bash
toggl - toggl cli

Usage:
  toggl-cli [command]

Available Commands:
  completion   Generate the autocompletion script for the specified shell
  delete-saved Delete a saved timer
  help         Help about any command
  new-saved    Save a new time entry
  pause        Pause the current entry
  report       Generate a report of a your time tracked
  resume       Resume the paused time entry
  start        Start new time entry
  start-saved  Start new time entry from your saved timers
  status       Get the curent tracking status
  stop         Stop the current entry

Flags:
  -h, --help   help for toggl-cli

Use "toggl-cli [command] --help" for more information about a command.
```

It's quite self expeditionary on how to use it from here. For a few quick tips
though:
* Add toggl timer to tmux status
```tmux
    set-option -g status-interval 1
    set -g status-right "#(toggl-cli status)"
```
    * This will add your current timer status to tmux and it will update the status every 1 second
* For some scripts and ways to quickly start timers you can check [the examples](https://github.com/skykosiner/toggl-cli/tree/master/examples)
