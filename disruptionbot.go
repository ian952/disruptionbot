package main

import (
	"fmt"
	"github.com/jpfuentes2/go-env"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	env.ReadEnv("./.env")
	go bindPort()
	go keepAppAwake()
	fmt.Printf("Using auth token: '%s'\n", os.Getenv("AUTH_TOKEN"))
	ws, _ := slackConnect(os.Getenv("AUTH_TOKEN"))
	fmt.Println("ready to disrupt")

	// XXX: Does not work if it has repeat letters
	feridun := strings.Split("FERIDUN", "")

	chatHistory := make(map[string]string)

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		if m.Type == "message" {
			content := strings.ToUpper(m.Text)
			channel := m.Channel

			cur := find(feridun, chatHistory[channel])
			if cur == -1 && content == feridun[0] {
				chatHistory[channel] = feridun[0]
			} else if content == feridun[cur+1] {
				if cur+1 == len(feridun)-1 {
					m.Text = "DISRUPTIVE!!!"
					postMessage(ws, m)
					chatHistory[channel] = ""
				} else {
					chatHistory[channel] = feridun[cur+1]
				}
			} else {
				chatHistory[channel] = ""
			}
		}
	}
}

func bindPort() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func keepAppAwake() {
	for {
		time.Sleep(time.Minute * 25)
		http.Get("https://disruption-bot.herokuapp.com/")
	}
}

func find(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}
