package main

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gleich/lumber"
	"github.com/gorilla/websocket"
	"github.com/nathany/bobblehat/sense/screen"
	"github.com/nathany/bobblehat/sense/screen/color"
)

type event struct {
	channel string
	text    string
	Type    string
	user    string
}

type pixel struct {
	channels []string
	updated  time.Time
}

func main() {
	lumber.Info("Booted up!")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "wss", Host: "joebunyan.haas.hackclub.com", Path: "stream"}
	lumber.Info("Connecting to", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		lumber.Fatal(err, "Failed to connect to", u.String())
	}

	fb := screen.NewFrameBuffer()
	err = screen.Clear()
	if err != nil {
		lumber.Fatal(err, "Failed to clear screen")
	}

	for {
		// Getting message
		_, message, err := c.ReadMessage()
		if err != nil {
			lumber.Error(err, "Failed to read message")
		}
		lumber.Success("Got message", string(message))

		// Parsing JSON from event
		var data event
		err = json.Unmarshal(message, &data)
		if err != nil {
			lumber.Error(err, "Failed to unmarshal", message)
		}

		// Setting pixel
		x := rand.Intn(8)
		y := rand.Intn(8)
		if data.Type == "message" {
			fb.SetPixel(x, y, color.White)
		} else if data.Type == "ping" {
			continue
		} else {
			fb.SetPixel(x, y, color.Green)
		}

		// Updating screen
		err = screen.Draw(fb)
		if err != nil {
			lumber.Error(err, "Failed to update screen")
		}

		// Having the pixel not update for 100 milliseconds
		go func() {
			time.Sleep(300 * time.Millisecond)
			fb.SetPixel(x, y, color.Black)
			err = screen.Draw(fb)
			if err != nil {
				lumber.Error(err, "Failed to set", x, ",", y, "coordinates")
			}
		}()
	}
}
