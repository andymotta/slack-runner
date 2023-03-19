# slack-runner
A SlackBot that runs simple scripts/commands and streams stdout to Slack
* Can also scheudle text and run shell commands

## Setup Notes

Environment variables:
Make sure you have set the `SLACK_APP_TOKEN` and `SLACK_BOT_TOKEN` environment variables in `docker-compose.yml`

Verify your Slack App settings:
Go to https://api.slack.com/apps and select your app.

Under "OAuth & Permissions", make sure you have the following bot token scopes added:

    app_mentions:read
    chat:write
    chat:write.public
    commands
    channels:history
    users:read

Under "Event Subscriptions", enable events, and subscribe to the app_mention event under "Subscribe to bot events".

Make sure "Socket Mode" is enabled.

## Run
```bash
docker-compose up --build
```