package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/slack-go/slack/slackevents"
)

var connections map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)

// Prevents two messages from being broadcast simultaneously
var mutex = sync.Mutex{}

func main() {
	go func() {
		for {
			for v := range connections {
				mutex.Lock()

				v.WriteJSON(map[string]interface{}{
					"type": "ping",
				})

				mutex.Unlock()
			}

			time.Sleep(5 * time.Second)
		}
	}()

	r := gin.Default()

	r.POST("/slack/events", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
			c.String(500, "sadge")
			return
		}

		event, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			log.Println(err)
			c.String(500, "sadge")
			return
		}

		if event.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse

			err := json.Unmarshal(body, &r)
			if err != nil {
				log.Println(err)
				c.String(500, "sadge")
				return
			}

			c.String(200, r.Challenge)
		} else if event.Type == slackevents.CallbackEvent {
			innerEvent := event.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				if ev.SubType == "" {
					for v := range connections {
						mutex.Lock()

						v.WriteJSON(map[string]interface{}{
							"type":    "message",
							"channel": ev.Channel,
						})

						mutex.Unlock()
					}
				}

			case *slackevents.ReactionAddedEvent:
				for v := range connections {
					mutex.Lock()

					v.WriteJSON(map[string]interface{}{
						"type":     "reaction",
						"channel":  ev.Item.Channel,
						"reaction": ev.Reaction,
					})

					mutex.Unlock()
				}
			}
		}
	})

	r.GET("/stream", func(c *gin.Context) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			c.String(500, "sadge")
			return
		}

		connections[ws] = true

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				delete(connections, ws)
				log.Println("Client disconnected")
				break
			}
		}
	})

	r.Run(":3000")
}
