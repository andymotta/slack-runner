package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleCmd(ev *slackevents.AppMentionEvent, api *slack.Client, client *socketmode.Client, command []string) {
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
}
