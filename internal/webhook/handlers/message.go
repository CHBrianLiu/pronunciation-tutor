package handlers

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

// Handle processes message webhook events.
func handleTextMessageEvent(r MessageReplier, e *linebot.Event) bool {
	msg, ok := e.Message.(*linebot.TextMessage)
	if !ok {
		log.Printf("%s message is not supposed to be handled by this handler.\n", msg.Type())
	}
	replyMessage := linebot.NewTextMessage(msg.Text)
	if _, err := r.ReplyMessage(e.ReplyToken, replyMessage).Do(); err != nil {
		log.Printf("error replying to message: %v\n", err)
		return false
	}
	return true
}
