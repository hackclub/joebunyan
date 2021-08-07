package main

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"os"
	"os/signal"

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

func main() {
	lumber.Info("Booted up!")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: "joebunyan.ngrok.io", Path: "stream"}
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
		_, message, err := c.ReadMessage()
		if err != nil {
			lumber.Error(err, "Failed to read message")
		}

		lumber.Success("Got message")

		var data event
		err = json.Unmarshal(message, &data)
		if err != nil {
			lumber.Error(err, "Failed to unmarshal", message)
		}

		if data.Type == "message" {
			fb.SetPixel(rand.Intn(8), rand.Intn(8), color.White)
		} else {
			fb.SetPixel(rand.Intn(8), rand.Intn(8), color.Green)
		}

		err = screen.Draw(fb)
		if err != nil {
			lumber.Error(err, "Failed to update screen")
		}
	}
}
