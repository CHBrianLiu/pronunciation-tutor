package main

import (
	"github.com/CHBrianLiu/pronunciation-tutor/internal/env"
	"github.com/CHBrianLiu/pronunciation-tutor/internal/webhook"
	"github.com/CHBrianLiu/pronunciation-tutor/internal/webhook/handlers"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
)

func main() {
	c, err := linebot.New(env.LineChannelSecret.GetValue(), env.LineChannelAccessToken.GetValue())
	if err != nil {
		log.Fatalf("Failed to create a line bot client: %v\n", err)
	}
	httpHandler := webhook.NewHandler(c, handlers.NewManager(c))
	http.Handle("/webhook", httpHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
