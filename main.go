package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"log"
	"net/http"
	"os"
)
var api = slack.New("")
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/gs", gsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "This is index - see gs!")
}

func gsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("gs hit")
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err)
	}
	body := buf.String()
	log.Println(body)
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: "me06Zd6Q1gYvyvgiJbj16ETU"}))
	if err != nil {
		log.Println("bad message?")
		w.WriteHeader(http.StatusInternalServerError)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		_, err = w.Write([]byte(r.Challenge))
		if err != nil {
			log.Println(err)
		}
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			_, _, err = api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			if err != nil {
				log.Println(err)
			}
		}
	}
}
