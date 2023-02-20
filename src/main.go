package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func main() {
	appToken := os.Getenv("SLACK_APP_TOKEN")
	if appToken == "" {

	}

	if !strings.HasPrefix(appToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set.\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	api := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	go func() {
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

					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						text := ev.Text
						command := strings.Fields(text)
						command = command[1:]
						switch firstWord := command[0]; firstWord {
						case "help":
							_, _, err := api.PostMessage(ev.Channel, slack.MsgOptionText("`cmd` to access <list_of_commands> or `<scriptname> <arguments>`", false))
							if err != nil {
								fmt.Printf("failed posting message: %v", err)
							}
						case "cmd": //only allowed commands here!
							cmd := exec.Command(command[1], command[2:]...)
							stdout, err := cmd.StdoutPipe()
							if err != nil {
								fmt.Println(err)
							}

							err = cmd.Start()
							fmt.Println("The command is running")
							if err != nil {
								fmt.Println(err)
							}

							channel, timestamp, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Running the command: `"+command[1]+"` with supplied arguments...", false))
							if err != nil {
								fmt.Printf("failed posting message: %v", err)
							}

							// print the output of the subprocess
							scanner := bufio.NewScanner(stdout)
							var str strings.Builder
							for scanner.Scan() {
								m := scanner.Text()
								str.WriteString(m + "\n")
								api.UpdateMessage(channel, timestamp, slack.MsgOptionText("```\n"+str.String()+"```", false))
								fmt.Println()
								if str.Len() > 3999 { // Break up text because of Slack limits
									str.Reset()
								}
							}
							cmd.Wait()
						case "schedule":
							const (
								layout = "2006-1-2 03:04PM"
							)
							loc, _ := time.LoadLocation("America/Los_Angeles")
							t, _ := time.ParseInLocation(layout, command[3]+" "+command[4], loc)
							api.ScheduleMessage(ev.Channel, strconv.FormatInt(t.Unix(), 10), slack.MsgOptionText("@glados scale "+command[1]+" "+command[2], false))
							api.PostMessage(ev.Channel, slack.MsgOptionText("Scheduled scaling event for "+command[1]+" to "+command[2]+" at: "+command[3]+" "+command[4], false))
						default:
							err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
								if err != nil {
									return err
								}
								script := info.Name()
								extension := filepath.Ext(script)
								basename := strings.TrimSuffix(script, extension)
								fmt.Println(basename)
								if basename == command[0] {
									switch extension {
									case ".sh":
										command = append([]string{"bash"}, command...)
									case ".py":
										command = append([]string{"python3"}, command...)
									case ".js":
										command = append([]string{"node"}, command...)
									case ".php":
										command = append([]string{"php"}, command...)
									default:
										log.Println("Unsupported extension, please see Dockerfile")
									}
									command[1] = script
									cmd := exec.Command(command[0], command[1:]...)
									stdout, err := cmd.StdoutPipe()
									if err != nil {
										log.Fatal(err)
									}
									err = cmd.Start()
									fmt.Println("The command is running")
									if err != nil {
										fmt.Println(err)
									}

									channel, timestamp, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Found script: "+script, false))
									if err != nil {
										fmt.Printf("failed posting message: %v", err)
									}

									// print the output of the subprocess
									scanner := bufio.NewScanner(stdout)
									var str strings.Builder
									for scanner.Scan() {
										m := scanner.Text()
										str.WriteString(m + "\n")
										api.UpdateMessage(channel, timestamp, slack.MsgOptionText("```\n"+str.String()+"```", false))
										if str.Len() > 3999 { // Break up text because of Slack limits
											str.Reset()
										}
									}
									cmd.Wait()
								} else {
									log.Println("Not running " + script) // all this does is print this message for every other script in the loop
								}
								return nil
							})
							if err != nil {
								log.Println(err)
							}
						}
					}
				default:
					client.Debugf("unsupported Events API event received")
				}
			}
		}
	}()

	client.Run()
}
