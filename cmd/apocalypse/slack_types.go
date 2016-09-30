package main

// SlackTextMessage defines the structure for posting a text message to a Slack channel
type SlackTextMessage struct {
	ResponseType string                `json:"response_type"`
	Text         string                `json:"text"`
	Attachments  []SlackTextAttachment `json:"attachments"`
}

// SlackTextAttachment defines the structure for attaching text to Slack messages
type SlackTextAttachment struct {
	Text string `json:"text"`
}

// SlackOAuthResponse defines the structure of a response from Slack after an OAuth handshake
type SlackOAuthResponse struct {
	Ok          bool   `json:"ok"`
	Warning     string `json:"warning"`
	Scope       string `json:"scope"`
	TeamName    string `json:"team_name"`
	TeamID      string `json:"team_id"`
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`

	IncomingWebhook struct {
		ChannelName      string `json:"channel"`
		ChannelID        string `json:"channel_id"`
		ConfigurationURL string `json:"configuration_url"`
		URL              string `json:"url"`
	} `json:"incoming_webhook"`

	Bot struct {
		UserID      string `json:"bot_user_id"`
		AccessToken string `json:"bot_access_token"`
		Warning     string `json:"warning"`
	} `json:"bot"`
}
