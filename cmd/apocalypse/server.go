package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// SlackMessage is a text message to send to a Slack channel
type SlackMessage struct {
	url       string
	message   string
	quip      string
	logFields log.Fields
}

// Server handles polling for changes and reporting to the Slack channels on change.
type Server struct {
	clientID     string // publicly-available Slack ID of this client
	clientSecret string // top-secret password with Slack for our clientID
	currentValue float32
	tokens       map[string]*Account // a map of tokens -> all info we have about an integration. Stored as JSON for our DB
	mutex        sync.Mutex
	dataFilePath string            // for now, the database is just a JSON dump of our 'tokens' map
	outChan      chan SlackMessage // queue of messages to be delivered to Slack channels
	quitChan     chan interface{}  // quit channel - closed when we need to wrap up and exit
	waitGroup    sync.WaitGroup    // used along with quitChan to keep track of pending work
}

// NewServer returns a new Server
func NewServer(clientID string, clientSecret string, dataFilePath string) (*Server, error) {
	tokenMap := make(map[string]*Account)

	// load the data if found
	tokenData, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		// allow this error
		log.WithFields(log.Fields{
			"area": "db",
		}).Warnf("Could not load data file on start-up: %s", err)
	} else {
		tm := make(map[string]*Account)
		if err := json.Unmarshal(tokenData, &tm); err == nil {
			tokenMap = tm
		}
	}

	return &Server{
		clientID:     clientID,
		clientSecret: clientSecret,
		tokens:       tokenMap,
		mutex:        sync.Mutex{},
		dataFilePath: dataFilePath,
		outChan:      make(chan SlackMessage, 10000),
		quitChan:     make(chan interface{}),
		waitGroup:    sync.WaitGroup{},
	}, nil
}

// save the token data - write lock should already be held
func (s *Server) saveTokenData() error {
	jsonData, err := json.Marshal(s.tokens)
	if err != nil {
		return fmt.Errorf("Error marshalling token data: %s", err)
	}

	backupFile := fmt.Sprintf("%s.%d", s.dataFilePath, time.Now().Unix())
	if err := copyFileContents(s.dataFilePath, backupFile); err != nil {
		// allow this error
		log.WithFields(log.Fields{
			"area": "db",
		}).Errorf("Could not store data file: %s", err)
	}

	err = ioutil.WriteFile(s.dataFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing data to file: %s", err)
	}
	return nil
}

// Run starts the service.
func (s *Server) Run() {
	rand.Seed(time.Now().UTC().UnixNano())

	// outgoing sender workers
	go func() {
		for slackMessage := range s.outChan {
			func() {
				defer s.waitGroup.Done()
				log.WithFields(slackMessage.logFields).Debugf("Sending message to channel")

				// retry loop
				attemptCount := 0
				for {
					attemptCount++
					if err := s.sendTextMessage(slackMessage.url, slackMessage.message, slackMessage.quip); err != nil {
						log.WithFields(slackMessage.logFields).Errorf("Error sending text message - retry attempt #%d/3: %s", attemptCount, err)
						if attemptCount >= 3 {
							break
						}
					} else {
						log.WithFields(slackMessage.logFields).Infof("Sent message to channel")
						break
					}
					time.Sleep(1 * time.Second)
				}
			}()
		}
	}()

	// 538 polling loop
	for {
		func() {
			s.waitGroup.Add(1)
			defer s.waitGroup.Done()

			// check for quit
			select {
			case <-s.quitChan:
				log.Infof("Received QUIT message from Run() - quitting")
				return
			default:
			}

			trumpChance, err := fetchTrumpChance()
			if err != nil {
				log.WithFields(log.Fields{
					"area": "fetch",
				}).Errorf("Error fetching data from 538: %s", err)
				return
			}
			log.WithFields(log.Fields{
				"area":  "data",
				"value": trumpChance,
			}).Debugf("Trump's chance fetched")

			s.mutex.Lock()
			defer s.mutex.Unlock()

			s.currentValue = trumpChance
			needToSave := false

			// loop through each team to see if there's a change
			for teamID := range s.tokens {
				team := s.tokens[teamID]
				if team.ReportedTrumpChance == trumpChance {
					log.WithFields(log.Fields{
						"area":     "data",
						"teamID":   team.TeamID,
						"teamName": team.TeamName,
						"value":    trumpChance,
					}).Debugf("Trump's chance hasn't changed for team")
					continue
				}

				// build the +/- context string
				contextStr := ""
				if team.ReportedTrumpChance > 0 {
					percentageChange := trumpChance - team.ReportedTrumpChance
					contextStr = fmt.Sprintf(" (%+.1f%%)", percentageChange)
				}
				msg := fmt.Sprintf("Chance of a Trump apocalypse: %.1f%%%s https://projects.fivethirtyeight.com/2016-election-forecast", trumpChance, contextStr)
				quip := randomQuip()
				logFields := log.Fields{
					"area":        "slack",
					"teamID":      team.TeamID,
					"teamName":    team.TeamName,
					"channelID":   team.IncomingWebhook.ChannelID,
					"channelName": team.IncomingWebhook.ChannelName,
					"message":     msg,
					"quip":        quip,
					"value":       trumpChance,
				}

				// queue up the outgoing message
				s.waitGroup.Add(1)
				s.outChan <- SlackMessage{
					url:       s.tokens[teamID].IncomingWebhook.URL,
					message:   msg,
					quip:      quip,
					logFields: logFields,
				}

				// for simplicity, assume the message does get sent, and update the database now
				needToSave = true
				team.ReportedTrumpChance = trumpChance
			}

			if needToSave {
				log.WithFields(log.Fields{
					"area": "db",
				}).Infof("Saving token data")

				if err := s.saveTokenData(); err != nil {
					log.WithFields(log.Fields{
						"area": "db",
					}).Errorf("Error saving token data: %s", err)
				}
			}
		}()

		time.Sleep(5 * time.Minute)
	}
}

// Stop running
func (s *Server) Stop() {
	s.waitGroup.Wait()
}

// send a Slack text message to a team's channel
func (s *Server) sendTextMessage(url string, body string, quip string) error {
	msg := SlackTextMessage{
		ResponseType: "in_channel",
		Text:         body,
		Attachments: []SlackTextAttachment{
			SlackTextAttachment{
				Text: quip,
			},
		},
	}
	respBytes, err := postJSON(url, msg)
	if err != nil {
		return fmt.Errorf("Error posting JSON to %s: %s", url, err)
	}
	log.Debugf("Sent text message - response: %s", string(respBytes))

	return nil
}

func (s *Server) handleTrump(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.WithFields(log.Fields{
			"request": "/trump",
		}).Errorf("Error parsing form: %s", err)
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	token := r.PostFormValue("token")
	team := r.PostFormValue("team_id")
	teamDomain := r.PostFormValue("team_domain")
	channelID := r.PostFormValue("channel_id")
	channelName := r.PostFormValue("channel_name")
	userID := r.PostFormValue("user_id")
	userName := r.PostFormValue("user_name")
	command := r.PostFormValue("command")
	text := r.PostFormValue("text")
	responseURL := r.PostFormValue("response_url")

	logFields := log.Fields{
		"request":      "/trump",
		"token":        token,
		"team":         team,
		"team_domain":  teamDomain,
		"channel_id":   channelID,
		"channel_name": channelName,
		"user_id":      userID,
		"user_name":    userName,
		"command":      command,
		"text":         text,
		"response_url": responseURL,
	}

	s.mutex.Lock()
	currentValue := s.currentValue
	s.mutex.Unlock()

	log.WithFields(logFields).Info("Received /trump request")

	// respond immediately with "in_channel" to tell Slack to show the original /trump command
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"response_type": "in_channel"}`)); err != nil {
		log.WithFields(logFields).Errorf("Error writing response_type:in_channel JSON: %s", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// send the response in a separate request to avoid scrolling issues in Slack
	s.waitGroup.Add(1)
	go func() {
		time.Sleep(500 * time.Millisecond)
		s.outChan <- SlackMessage{
			url:       responseURL,
			message:   fmt.Sprintf("Chance of a Trump apocalypse: %.1f%% https://projects.fivethirtyeight.com/2016-election-forecast", currentValue),
			quip:      randomQuip(),
			logFields: logFields,
		}
	}()
}

// handle incoming OAuth requests
func (s *Server) handleOAuth(w http.ResponseWriter, r *http.Request) {
	select {
	case <-s.quitChan:
		w.Write([]byte("Please try again later"))
		return
	default:
	}

	if errorReason := r.URL.Query().Get("error"); errorReason == "access_denied" {
		// user clicked "Cancel"
		w.Write([]byte("Maybe next time!"))
		return
	}

	oauthCode := r.URL.Query().Get("code")
	if oauthCode == "" {
		log.Errorf("Error: handleOAuth missing 'code'")
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	oauthURL := "https://slack.com/api/oauth.access"
	params := url.Values{}
	params.Add("client_id", s.clientID)
	params.Add("client_secret", s.clientSecret)
	params.Add("code", oauthCode)
	requestStr := params.Encode()

	logFields := log.Fields{
		"area": "oath",
	}

	respBytes, err := postRequest(oauthURL, requestStr)
	log.WithFields(logFields).Debugf("Posting to %s", oauthURL)
	if err != nil {
		log.WithFields(logFields).Errorf("Error posting to %s: %s", oauthURL, err)
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	log.WithFields(logFields).Debugf("oauth response: %s", string(respBytes))

	oauthResponse := Account{}
	err = json.Unmarshal(respBytes, &oauthResponse)
	if err != nil {
		log.WithFields(logFields).Errorf("Error unmarshalling Account: %s", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	if oauthResponse.AccessToken == "" {
		log.WithFields(logFields).Errorf("Empty AccessToken")
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	if oauthResponse.TeamID == "" {
		log.WithFields(logFields).Errorf("Empty TeamID")
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	s.mutex.Lock()
	s.tokens[oauthResponse.TeamID] = &oauthResponse
	if err := s.saveTokenData(); err != nil {
		// allow
		log.WithFields(logFields).Errorf("Error saving token data: %s", err)
	}
	s.mutex.Unlock()

	log.WithFields(logFields).Infof("Successfully received access token: %s with method %s", oauthResponse.AccessToken, oauthResponse.Scope)

	http.Redirect(w, r, oauthResponse.IncomingWebhook.ConfigurationURL, http.StatusTemporaryRedirect)
}
