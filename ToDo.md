## Temp delete scheduled event
```bash
# List events to grab the one you want
curl -H "Authorization: Bearer xoxb-<token>" https://slack.com/api/chat.scheduledMessages.list
# Example response
{"ok":true,"scheduled_messages":[{"id":"Q025W0XUFTL","channel_id":"C013YHQ7UA2","post_at":1624575600,"date_created":1624402134,"text":"Test message!"}],"response_metadata":{"next_cursor":""}}

# Delete event by ID output above
curl -XPOST -H "Authorization: Bearer xoxb-<token>" --data "channel=C013YHQ7UA2&scheduled_message_id=Q025W0XUFTL" https://slack.com/api/chat.deleteScheduledMessage
# response
{"ok":true}
```