package handlers_test

import (
	"encoding/json"
	"github.com/CHBrianLiu/pronunciation-tutor/internal/webhook/handlers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type testReplyMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

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
	expectedReplyMsg := testReplyMessage{
		Type: "text",
		Text: msgReceived,
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				// Test method and path
				assert.Equal(t, "POST", req.Method)
				assert.Equal(t, "/v2/bot/message/reply", req.URL.String())
				// Unmarshal the request data to check the content
				dec := json.NewDecoder(req.Body)
				data := &struct {
					ReplyToken           string             `json:"replyToken"`
					Messages             []testReplyMessage `json:"messages"`
					NotificationDisabled bool               `json:"notificationDisabled,omitempty"`
				}{}
				err := dec.Decode(data)
				assert.NoError(t, err)
				// Test the request data
				assert.Equal(t, replyToken, data.ReplyToken)
				assert.Len(t, data.Messages, 1)
				assert.Equal(t, expectedReplyMsg, data.Messages[0])
			}))
	httpClient := server.Client()

	// Close the server when test finishes
	defer server.Close()

	c, err := linebot.New(
		"secret",
		"token",
		// Currently not sure the purpose of this option.
		linebot.WithHTTPClient(httpClient),
		// We need this option to send requests to our test server.
		linebot.WithEndpointBase(server.URL),
	)
	assert.NoError(t, err)

	handler, ok := handlers.NewManager(c).GetHandler(event)
	assert.True(t, ok, "Message event handler not found.")
	// Assert the Line API call in the test server handler function.
	handler(event)
}
