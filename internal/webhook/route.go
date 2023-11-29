package webhook

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"net/http"
)

type RequestParser interface {
	// ParseRequest should validate the request and parse it into events.
	ParseRequest(r *http.Request) ([]*linebot.Event, error)
}

type HandlerManager interface {
	// GetHandler returns the proper handler function by examining the event.
	// If no handler is registered for certain types of events, return (nil, false).
	GetHandler(*linebot.Event) (func(*linebot.Event), bool)
}

// Handler implements http.Handler and can serve Line webhook HTTP requests.
type Handler struct {
	parser  RequestParser
	manager HandlerManager
}

// ServeHTTP handles a Line webhook request.
// It creates goroutines for each event if corresponding event handlers exist.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	events, err := h.parser.ParseRequest(r)
	if err != nil {
		fmt.Printf("Failed to parse the request with error: %v\n", err)
		w.WriteHeader(400)
		return
	}
	for _, event := range events {
		handler, ok := h.manager.GetHandler(event)
		if !ok {
			fmt.Printf("No event handler registered for type %s, ignore the event\n", event.Type)
			continue
		}
		go handler(event)
	}
	w.WriteHeader(200)
}

// NewHandler constructs a new Handler and returns a pointer to it.
func NewHandler(p RequestParser, m HandlerManager) *Handler {
	return &Handler{
		parser:  p,
		manager: m,
	}
}
