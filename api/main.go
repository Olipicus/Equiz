package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"code.olipicus.com/equiz/api/equiz"
	"code.olipicus.com/equiz/api/line"
	"code.olipicus.com/equiz/config"
)

var appConfig *config.Table

func main() {
	port := flag.Int("port", 8080, "port to listen")
	configPath := flag.String("config", "./config.dev.json", "config path")
	flag.Parse()

	appConfig = config.LoadConfig(*configPath)
	service := equiz.New(*configPath)

	channelSecret, _ := appConfig.GetString("channel_secret")
	channelToken, _ := appConfig.GetString("channel_token")

	app, err := line.NewLineApp(channelSecret, channelToken, service)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", app.CallbackHandler)
	http.HandleFunc("/test", testHandler)

	fmt.Printf("Server run on: %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

}

func testHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(appConfig)

	if err != nil {
		w.Write([]byte("Get Config Error"))
	}
	w.Write(b)

}
