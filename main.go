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
var authToken = os.Getenv("SLACK_TOKEN")
var verifToken = os.Getenv("SLACK_VERIF_TOKEN")

var api = slack.New(authToken)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if authToken == "" {
		log.Println("bad auth token")
	}
	if verifToken == "" {
		log.Println("bad verif token")
	}

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
	messageEvent, err := slackevents.ParseActionEvent(body, slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verifToken}))
	if err != nil {
		log.Println("json bad..")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Println(messageEvent)

	if messageEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text")
		_, err = w.Write([]byte(r.Challenge))
		if err != nil {
			log.Println(err)
		}
	}
	_, _, err = api.PostMessage(messageEvent.Channel.ID, slack.MsgOptionText("Yes, hello.", false))
	if err != nil {
		log.Println(err)
	}
}
