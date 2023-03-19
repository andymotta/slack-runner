package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func handleHelp(ev *slackevents.AppMentionEvent, api *slack.Client, client *socketmode.Client, command []string) {
	scriptsPath := "../scripts"
	var scripts []string

	files, err := ioutil.ReadDir(scriptsPath)
	if err != nil {
		fmt.Printf("Error reading the scripts directory: %v", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			script := file.Name()
			extension := filepath.Ext(script)
			basename := strings.TrimSuffix(script, extension)
			scripts = append(scripts, basename)
		}
	}

	scriptList := strings.Join(scripts, ", ")
	helpText := fmt.Sprintf("`cmd` run anything in the bot shell, `schedule` to schedule a message or one of following commands:\n```\n%s\n```", scriptList)

	_, _, err = api.PostMessage(ev.Channel, slack.MsgOptionText(helpText, false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}
