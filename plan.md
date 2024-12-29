# Commands
- [x] toggl status
    * Print the current status
- [x] toggl pause
    * Pause the current time tracking
- [x] toggl resume
    * Resume the current tracking
- [x] toggl stop
    * Stop the current tracking entry
- [x] toggl start-saved
    * Start time tracking with saved timers in your config
- [x] toggl start
    * List all your projects in general
        * Ask if you want to add a tag and then ask if you want to add a description
- [x] toggl new-saved
    * List your projects
        * Select a project
    * List tags
        * Select tag(s)
    * Ask if you want to add a description
- [x] toggl delete-saved
    * List all your saved timers then select which ones you want to delete
- [x] toggl report day|week|month|year
    * Get report information for what the arg is

# Config layout
```json
{
    "api_key": "xxx",
    "workspace_id": "xxx",
    "saved_timers": [
        {
            "name": "Coding Toggl Project",
            // ID of the project you want to time track
            "project_id": "xxx",
            // ID of each tag you want to add
            "tags": [
                "xxx",
                "xxxx"
            ],
            // Leave blank for no description
            "description": ""
        }
    ]
}
```
