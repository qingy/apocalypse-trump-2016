package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// FetchTrumpChance fetches the chance that Trump will win the election
func fetchTrumpChance() (float32, error) {
	doc, err := goquery.NewDocument("http://projects.fivethirtyeight.com/2016-election-forecast")
	if err != nil {
		return 0, fmt.Errorf("Error fetching document: %s", err)
	}

	percentStr := doc.Find("[data-card-id='US-winprob-sentence'] .candidate-val.winprob[data-key='winprob'][data-party='R']").Text()
	if !strings.HasSuffix(percentStr, "%") {
		return 0, fmt.Errorf("Error finding percentage")
	}

	percentStr = strings.TrimSuffix(percentStr, "%")
	val, err := strconv.ParseFloat(percentStr, 32)
	return float32(val), err
}
