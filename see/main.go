package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gleich/lumber"
	"github.com/gorilla/websocket"
	"github.com/nathany/bobblehat/sense/screen"
	"github.com/nathany/bobblehat/sense/screen/color"
)

type event struct {
	Channel string
	Type    string
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
			lumber.Error(err, "Failed to unmarshal", string(message))
		}

		// Ignore pings
		if data.Type == "ping" {
			continue
		}

		// Setting pixel
		x := genLoc(data.Channel, true)
		y := genLoc(data.Channel, false)
		var pixelColor color.Color
		if data.Type == "message" {
			if len(os.Args) > 1 && os.Args[1] == "--christmas" {
				pixelColor = color.Red
			} else {
				pixelColor = color.White
			}
		} else {
			pixelColor = color.Green
		}
		fb.SetPixel(x, y, pixelColor)

		// Updating screen
		err = screen.Draw(fb)
		if err != nil {
			lumber.Error(err, "Failed to update screen")
		}

		// Having the pixel not update for 100 milliseconds
		go func() {
			time.Sleep(200 * time.Millisecond)
			fb.SetPixel(x, y, color.Black)
			err = screen.Draw(fb)
			if err != nil {
				lumber.Error(err, "Failed to set", x, ",", y, "coordinates")
			}
		}()
	}
}

// Take a string and if it should be the first four or last four characters, get the 32-bit FNV-1a hash encoding of the string
// then convert that long number to a single digit and cap it at 7
func genLoc(str string, first bool) int {
	fourChars := str[:4]
	if !first {
		fourChars = str[:len(str)-4]
	}
	h := fnv.New32a()
	h.Write([]byte(fourChars))
	loc, err := strconv.Atoi(fmt.Sprint(h.Sum32())[:1])
	if err != nil {
		lumber.Error(err, "Failed to convert", fourChars, "to coordinate")
	}
	if loc > 7 {
		return 7
	}
	return loc
}
