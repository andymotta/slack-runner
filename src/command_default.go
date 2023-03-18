package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleDefault(ev *slackevents.AppMentionEvent, api *slack.Client, client *socketmode.Client, command []string) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		script := info.Name()
		extension := filepath.Ext(script)
		basename := strings.TrimSuffix(script, extension)

		if basename == command[0] {
			executeScript(api, ev, script, extension, command)
		} else {
			log.Println("Not running " + script)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

func executeScript(api *slack.Client, ev *slackevents.AppMentionEvent, script, extension string, command []string) {
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
		return
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
}
