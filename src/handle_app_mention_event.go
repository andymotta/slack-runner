package main

import (
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleAppMentionEvent(ev *slackevents.AppMentionEvent, api *slack.Client, client *socketmode.Client) {
	text := ev.Text
	command := strings.Fields(text)
	if len(command) < 2 {
		api.PostMessage(ev.Channel, slack.MsgOptionText("No command provided. Type `@bot help` for a list of available commands.", false))
		return
	}
	command = command[1:]
	switch firstWord := command[0]; firstWord {
	case "help":
		handleHelp(ev, api, client, command)
	case "cmd":
		handleCmd(ev, api, client, command)
	case "schedule":
		handleSchedule(ev, api, client, command)
	default:
		handleDefault(ev, api, client, command)
	}
}
