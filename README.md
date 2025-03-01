# Hnotify
Hnotify is a Go-based tool designed to monitor changes in bug bounty programs (on HackerOne) and send notifications to Discord when new programs or assets are added. It fetches data from the arkadiyt/bounty-targets-data repository and compares it with a local copy to detect changes.

### Requirements
- HNOTIFY_DISCORD_WEBHOOK_URL: This is the only required environment variable. It specifies the Discord webhook URL where notifications will be sent.
- Optional Variables: For additional customization, you can set other environment variables. Check the main.go source code for details.
