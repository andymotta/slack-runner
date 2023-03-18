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
