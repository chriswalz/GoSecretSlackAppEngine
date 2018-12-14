package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/jsonapi"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
)

func main() {
	log.Println(runtime.GOOS)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
	_, err := fmt.Fprint(w, "GoSecret Slack engine is running. Have slack command call /gs endpoint!")
	if err != nil {
		log.Print(err)
	}
}

func gsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("gs hit")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)

	}

	// receive message from slack server
	type PayloadReq struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Message  string `json:"message"`
				Password string `json:"password"`
			} `json:"attributes"`
		} `json:"data"`
	}
	payloadReq := PayloadReq{}
	payloadReq.Data.Type = "Secrets"
	payloadReq.Data.Attributes.Message = r.FormValue("text")
	jsonBody, err := json.Marshal(&payloadReq)
	if err != nil {
		log.Println(err)
	}

	// send message to GoSecret.io then receive response from GoSecret.io
	resp, err := http.Post("https://www.gosecret.io/api/create", jsonapi.MediaType, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
	}
	payloadResp := PayloadReq{}
	bodyB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(bodyB, &payloadResp)
	if err != nil {
		log.Println(err)
	}

	jsonResp, _ := json.Marshal(struct {
		Type string `json:"response_type"`
		Text string `json:"text"`
	}{
		Type: "in_channel",
		Text: fmt.Sprintf("_GoSecret message made with_ \"/gs [msg]\"\n%s", "https://www.gosecret.io"+payloadResp.Data.Attributes.Message),
	})
	// send to a response url (instead of using ResponseWriter) for a "delayed response" to prevent command text from being echoed in Slack
	_, _ = http.Post(r.FormValue("response_url"), jsonapi.MediaType, bytes.NewBuffer(jsonResp))
	if err != nil {
		log.Println(err)
	}
}
