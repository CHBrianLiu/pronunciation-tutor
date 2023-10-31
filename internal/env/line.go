package env

var (
	// LineChannelAccessToken is used to send messages to the Line messaging API.
	LineChannelAccessToken = newEnv("LINE_CHANNEL_ACCESS_TOKEN")
	// LineChannelSecret is used to verify the webhook messages.
	LineChannelSecret = newEnv("LINE_CHANNEL_SECRET")
)
