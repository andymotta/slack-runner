package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleSchedule(ev *slackevents.AppMentionEvent, api *slack.Client, client *socketmode.Client, command []string) {
	if len(command) < 2 {
		api.PostMessage(ev.Channel, slack.MsgOptionText("Invalid schedule command. Usage: `schedule [list|delete|create] [args]`", false))
		return
	}

	action := command[1]

	switch action {
	case "list":
		listScheduledEvents(api, ev.Channel)
	case "delete":
		if len(command) < 3 {
			api.PostMessage(ev.Channel, slack.MsgOptionText("Please provide the ID of the scheduled message to delete. Usage: `schedule delete [message_id]`", false))
			return
		}
		deleteScheduledEvent(api, ev.Channel, command[2])
	case "create":
		createScheduledEvent(api, ev.Channel, command[2:])
	default:
		api.PostMessage(ev.Channel, slack.MsgOptionText("Invalid action. Usage: `schedule [list|delete|create] [args]`", false))
	}
}

func listScheduledEvents(api *slack.Client, channel string) {
	scheduledMessages, _, err := api.GetScheduledMessages(&slack.GetScheduledMessagesParameters{Channel: channel})
	if err != nil {
		log.Printf("Error getting scheduled messages: %v", err)
		api.PostMessage(channel, slack.MsgOptionText(fmt.Sprintf("Error getting scheduled messages: %v", err), false))
		return
	}

	if len(scheduledMessages) == 0 {
		api.PostMessage(channel, slack.MsgOptionText("No scheduled messages found", false))
		return
	}

	var message strings.Builder
	message.WriteString("Scheduled messages:\n")
	for _, scheduledMessage := range scheduledMessages {
		timestamp := int64(scheduledMessage.PostAt)
		timeStr := time.Unix(timestamp, 0).Format(time.RFC1123)
		message.WriteString(fmt.Sprintf("- Message ID: %s, Scheduled time: %s\n", scheduledMessage.ID, timeStr))
	}
	api.PostMessage(channel, slack.MsgOptionText(message.String(), false))
}

func deleteScheduledEvent(api *slack.Client, channel string, messageID string) {
	params := &slack.DeleteScheduledMessageParameters{
		Channel:            channel,
		ScheduledMessageID: messageID,
	}
	successful, err := api.DeleteScheduledMessage(params)
	if err != nil {
		log.Printf("Error deleting scheduled message: %v", err)
		api.PostMessage(channel, slack.MsgOptionText(fmt.Sprintf("Error deleting scheduled message: %v", err), false))
		return
	}
	if !successful {
		api.PostMessage(channel, slack.MsgOptionText(fmt.Sprintf("Scheduled message %s not found", messageID), false))
		return
	}
	api.PostMessage(channel, slack.MsgOptionText(fmt.Sprintf("Scheduled message %s deleted", messageID), false))
}

func createScheduledEvent(api *slack.Client, channel string, command []string) {
	if len(command) < 3 {
		api.PostMessage(channel, slack.MsgOptionText("Invalid create schedule command. Usage: `schedule create [date] [time] [message]`", false))
		return
	}

	const (
		layout = "2006-1-2 03:04PM"
	)
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Printf("Error loading location: %v", err)
		api.PostMessage(channel, slack.MsgOptionText("Error loading location. Please check the timezone settings.", false))
		return
	}

	t, err := time.ParseInLocation(layout, command[0]+" "+command[1], loc)
	if err != nil {
		api.PostMessage(channel, slack.MsgOptionText("Invalid date or time format. Please use the format `YYYY-MM-DD hh:mmAM/PM`", false))
		return
	}

	message := strings.Join(command[2:], " ")

	_, _, err = api.ScheduleMessage(channel, strconv.FormatInt(t.Unix(), 10), slack.MsgOptionText(message, false))
	if err != nil {
		log.Printf("Error scheduling message: %v", err)
		api.PostMessage(channel, slack.MsgOptionText(fmt.Sprintf("Error scheduling message: %v", err), false))
		return
	}
	api.PostMessage(channel, slack.MsgOptionText("Scheduled message for "+command[0]+" "+command[1]+": "+message, false))
}
