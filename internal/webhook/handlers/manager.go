package handlers

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

type eventHandler = func(*linebot.Client, *linebot.Event) bool

type HandlerCollection map[linebot.EventType]eventHandler

var handlers = HandlerCollection{
	linebot.EventTypeMessage: handleTextMessageEvent,
}

type Manager struct {
	handlers HandlerCollection
	replier  *linebot.Client
}

func (m *Manager) GetHandler(event *linebot.Event) (func(*linebot.Event), bool) {
	handler, ok := m.handlers[event.Type]
	if !ok {
		return nil, false
	}
	return func(event *linebot.Event) {
		if ok = handler(m.replier, event); !ok {
			log.Printf("Failed to handle webhook event %s\n", event.Type)
		}
	}, true
}

func NewManager(r *linebot.Client) *Manager {
	return &Manager{handlers, r}
}
