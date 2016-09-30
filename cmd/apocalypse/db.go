package main

// Account is a record we store in the DB, holding everything
// we need about a connected Slack account.
type Account struct {
	SlackOAuthResponse
	ReportedTrumpChance float32 `json:"reported_trump_stance"`
}
