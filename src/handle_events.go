package main

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleEvents(client *socketmode.Client, api *slack.Client) {
	for evt := range client.Events {
		switch evt.Type {
		case socketmode.EventTypeConnecting:
			fmt.Println("Connecting to Slack with Socket Mode...")
		case socketmode.EventTypeConnectionError:
			fmt.Println("Connection failed. Retrying later...")
		case socketmode.EventTypeConnected:
			fmt.Println("Connected to Slack with Socket Mode.")
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				fmt.Printf("Ignored %+v\n", evt)
			}

			fmt.Printf("Event received: %+v\n", eventsAPIEvent)

			client.Ack(*evt.Request)

			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				innerEvent := eventsAPIEvent.InnerEvent
				switch ev := innerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					handleAppMentionEvent(ev, api, client)
				case *slackevents.MessageEvent:
					handleMessageEvent(ev, api, client)
				}
			default:
				client.Debugf("unsupported Events API event received")
			}
		}
	}
}
