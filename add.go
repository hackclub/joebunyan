// This script will add the joebunyan slack bot to every public channel in the hackclub slack org
package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	client := &http.Client{}
	channels, err := getChannels(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("\n\nGot %v channels\n\n\n", len(channels))
}

func getChannels(client *http.Client) ([]string, error) {
	channels := []string{}
	cursor := ""
	for {
		url := "https://slack.com/api/conversations.list?exclude_archived=true&types=public_channel&cursor=" + url.QueryEscape(cursor)
		req, err := http.NewRequest("GET", url, strings.NewReader(""))
		if err != nil {
			return []string{}, err
		}

		req.Header.Add("Authorization", "Bearer "+os.Getenv("ADD_OAUTH_TOKEN"))

		res, err := client.Do(req)
		if err != nil {
			return []string{}, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return []string{}, err
		}

		var data struct {
			OK       bool
			Error    string
			Channels []struct {
				ID        string
				Is_shared bool `json:"is_shared"`
			}
			ResponseMetadata struct {
				NextCursor string `json:"next_cursor"`
			} `json:"response_metadata"`
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return []string{}, err
		}

		if !data.OK {
			if data.Error == "ratelimited" {
				log.Println("Currently ratelimited. Resting for one minute")
				time.Sleep(time.Minute)
				continue
			}
			return []string{}, errors.New("Data returned non OK from slack API: " + string(body))
		}

		if data.ResponseMetadata.NextCursor == "" {
			break
		}

		for _, channel := range data.Channels {
			if !channel.Is_shared {
				channels = append(channels, channel.ID)
			}
		}
		log.Printf("Got %v channels so far\n", len(channels))

		cursor = data.ResponseMetadata.NextCursor
	}
	return channels, nil
}
