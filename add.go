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
	token := "Bearer " + os.Getenv("ADD_OAUTH_TOKEN")

	channels, err := getChannels(client, token)
	if err != nil {
		log.Fatal(err)
	}

	err = joinChannels(client, channels, token)
	if err != nil {
		log.Fatal(err)
	}
}

func getChannels(client *http.Client, token string) (map[string]string, error) {
	channels := map[string]string{}
	cursor := ""
	for {
		req, err := http.NewRequest(
			"GET",
			"https://slack.com/api/conversations.list?exclude_archived=true&types=public_channel&cursor="+url.QueryEscape(
				cursor,
			),
			strings.NewReader(""),
		)
		if err != nil {
			return map[string]string{}, err
		}

		req.Header.Add("Authorization", token)

		res, err := client.Do(req)
		if err != nil {
			return map[string]string{}, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return map[string]string{}, err
		}

		var data struct {
			OK       bool
			Error    string
			Channels []struct {
				ID       string
				Name     string
				IsShared bool `json:"is_shared"`
			}
			ResponseMetadata struct {
				NextCursor string `json:"next_cursor"`
			} `json:"response_metadata"`
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return map[string]string{}, err
		}

		if !data.OK {
			if data.Error == "ratelimited" {
				log.Println("Currently ratelimited. Resting for one minute")
				time.Sleep(time.Minute)
				continue
			}
			return map[string]string{}, errors.New(
				"Data returned non OK from slack API: " + string(body),
			)
		}

		if data.ResponseMetadata.NextCursor == "" {
			break
		}

		for _, channel := range data.Channels {
			if !channel.IsShared {
				channels[channel.ID] = channel.Name
			}
		}
		log.Printf("Got %v channels so far\n", len(channels))

		cursor = data.ResponseMetadata.NextCursor
	}
	return channels, nil
}

func joinChannels(client *http.Client, channels map[string]string, token string) error {
	stageAdded := 0
	added := 0
	for id, name := range channels {
		if stageAdded >= 50 {
			log.Println("Sleeping for 1 minute to prevent rate limiting")
			time.Sleep(time.Minute)
			stageAdded = 0
		}

		req, err := http.NewRequest("POST", "https://slack.com/api/conversations.join?channel="+id,
			strings.NewReader(""),
		)
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", token)

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var data struct {
			OK    bool
			Error string
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}

		if data.Error != "" || !data.OK {
			return errors.New(data.Error)
		}
		added++
		stageAdded++
		log.Printf("Added to #%v - %v/%v\n", name, added, len(channels))
	}
	return nil
}
