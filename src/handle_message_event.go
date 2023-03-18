package main

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleMessageEvent(ev *slackevents.MessageEvent, api *slack.Client, client *socketmode.Client) {
	// Check if the message contains the word "hello" (case-insensitive)
	if strings.Contains(strings.ToLower(ev.Text), "hello") {
		// Reply with a "Hello!" message
		_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Hello!", false))
		if err != nil {
			fmt.Printf("Failed to post message: %v", err)
		}
	}
}
