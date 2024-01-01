package handlers_test

import (
	"github.com/CHBrianLiu/pronunciation-tutor/internal/webhook/handlers"
	"github.com/CHBrianLiu/pronunciation-tutor/mocks"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func TestMessageEventHandler_EchoTextMessage(t *testing.T) {
	replyToken := "fake-reply-token"
	msgReceived := "hello, world"
	event := &linebot.Event{
		Type: linebot.EventTypeMessage,
		Message: &linebot.TextMessage{
			Text: msgReceived,
		},
		ReplyToken: replyToken,
	}
	expectedReplyMessage := linebot.NewTextMessage(msgReceived)

	replier := mocks.NewMessageReplier(t)
	doer := mocks.NewReplyDoer(t)
	doer.On("Do").Return(&linebot.BasicResponse{RequestID: "request-id"}, nil)
	replier.On("ReplyMessage", replyToken, expectedReplyMessage).Return(doer)

	handler, ok := handlers.NewManager(replier).GetHandler(event)
	assert.True(t, ok, "Message event handler not found.")
	handler(event)

	replier.AssertCalled(t, "ReplyMessage", replyToken, expectedReplyMessage)
	doer.AssertCalled(t, "Do")
}
