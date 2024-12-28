# Commands
* toggl status
    * Print the current status
* toggl pause
    * Pause the current time tracking
* toggl resume
    * Resume the current tracking
* toggl stop
    * Stop the current tracking entry
* toggl start
    * List your saved timers in config
    * List all your projects in general
        * Ask if you want to add a tag and then ask if you want to add a description
* toggl new-saved
    * List your projects
        * Select a project
    * List tags
        * Select tag(s)
    * Ask if you want to add a description
* toggl delete-saved
    * List all your saved timers then select which ones you want to delete
* toggl report day|week|month|year
    * Get report information for what the arg is

# Config layout
```jsonc
{
    "api_key": "xxx",
    "saved_timers": [
        {
            "name": "Coding Toggl Project",
            // ID of the project you want to time track
            "project_id": "xxx",
            // ID of each tag you want to add
            "tags": [
                "xxx",
                "xxxx",
            ],
            // Leave blank for no description
            "description": ""
        }
    ]
}
```
