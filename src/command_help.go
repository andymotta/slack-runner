package main

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleHelp(ev *slackevents.AppMentionEvent, api *slack.Client, client *socketmode.Client, command []string) {
	_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("`cmd` to access <list_of_commands> or `<scriptname> <arguments>`", false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}
