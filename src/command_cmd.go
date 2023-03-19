package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleCmd(ev *slackevents.AppMentionEvent, api *slack.Client, client *socketmode.Client, command []string) error {
	cmd := exec.Command(command[1], command[2:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error creating StdoutPipe:", err)
		return err
	}

	err = cmd.Start()
	if err != nil {
		log.Println("Error starting command:", err)
		return err
	}
	fmt.Println("The command is running")

	channel, timestamp, err := api.PostMessage(ev.Channel, slack.MsgOptionText("Running the command: `"+command[1]+"` with supplied arguments...", false))
	if err != nil {
		log.Printf("failed posting message: %v", err)
		return err
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
	err = cmd.Wait()
	if err != nil {
		log.Println("Error waiting for command to finish:", err)
		return err
	}
	return nil
}
