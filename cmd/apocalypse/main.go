package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	var dataFilePath string
	var logLevel string
	var listenOn string
	var rootRedirectLocation string

	flag.StringVar(&dataFilePath, "data-file-path", "", "Location of the JSON DB file")
	flag.StringVar(&logLevel, "log-level", "info", "log level: debug, info, warning, error, fatal, panic")
	flag.StringVar(&listenOn, "listen", "", "<host>:<port> to listen on")
	flag.StringVar(&rootRedirectLocation, "root-redirect", "", "Where to redirect for /")

	flag.Usage = func() {
		fmt.Println("apocalypse2016 usage:")
		flag.PrintDefaults()
		fmt.Println("\nIn addition, the following environment variables are required:")
		fmt.Println("  CLIENT_ID\n    \tSlack client ID")
		fmt.Println("  CLIENT_SECRET\n    \tSlack client secret")
	}
	flag.Parse()

	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		fmt.Printf("Invalid log level (%s): %s\n\n", logLevel, err)
		flag.Usage()
		os.Exit(-1)
	}
	log.SetLevel(lvl)

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	if clientID == "" || clientSecret == "" || dataFilePath == "" || listenOn == "" {
		flag.Usage()
		os.Exit(-1)
	}

	server, err := NewServer(clientID, clientSecret, dataFilePath)
	if err != nil {
		fmt.Printf("Error instantiating Server: %s\n", err)
		os.Exit(-1)
	}

	// watch for ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			log.Infof("Received interrupt - waiting for all work to be done")
			server.Stop()
			log.Infof("Work is done - exiting")
			os.Exit(0)
		}
	}()

	go server.Run()

	// HTTP endpoints:
	if rootRedirectLocation != "" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, rootRedirectLocation, http.StatusTemporaryRedirect)
		})
	}
	http.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		server.handleOAuth(w, r)
	})
	http.HandleFunc("/trump", func(w http.ResponseWriter, r *http.Request) {
		server.handleTrump(w, r)
	})

	err = http.ListenAndServe(listenOn, nil)
	if err != nil {
		fmt.Printf("Error listening on %s: %s\n", listenOn, err)
		os.Exit(-1)
	}
}
