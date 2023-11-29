package webhook_test

import (
	"fmt"
	"github.com/CHBrianLiu/pronunciation-tutor/internal/webhook"
	"github.com/CHBrianLiu/pronunciation-tutor/mocks"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandler_ServeHTTP_ValidationError tests the scenario that the incoming request
// is not valid. In this situation, we respond with the status code 400 Bad Request.
func TestHandler_ServeHTTP_ValidationError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	res := httptest.NewRecorder()
	parser := mocks.NewRequestParser(t)
	parser.
		On("ParseRequest", req).
		Return(
			[]*linebot.Event{},
			fmt.Errorf("failed to parse the request into line webhook events"),
		)

	handler := webhook.NewHandler(parser, mocks.NewHandlerManager(t))

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

// TestHandler_ServeHTTP_DispatchEvents test the scenario that one join event and one leave event
// are sent to the webhook handler, and we only have join event handler registered.
// In this situation, we create a join-event-handler goroutine and ignore the other.
func TestHandler_ServeHTTP_DispatchEvents(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	joinEvent := linebot.Event{Type: linebot.EventTypeJoin}
	leaveEvent := linebot.Event{Type: linebot.EventTypeLeave}
	// parser returns two events directly.
	parser := mocks.NewRequestParser(t)
	parser.
		On("ParseRequest", req).
		Return(
			[]*linebot.Event{&leaveEvent, &joinEvent},
			nil,
		)
	// manager has only the join event handler registered.
	manager := mocks.NewHandlerManager(t)
	doneSig := make(chan linebot.EventType)
	joinEvHandler := func(e *linebot.Event) {
		doneSig <- e.Type
	}
	manager.
		On("GetHandler", &joinEvent).
		Return(joinEvHandler, true).
		Once()
	manager.
		On("GetHandler", &leaveEvent).
		Return(nil, false).
		Once()
	res := httptest.NewRecorder()

	handler := webhook.NewHandler(parser, manager)

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	select {
	case handled := <-doneSig:
		assert.Equal(t, linebot.EventTypeJoin, handled)
	}
}
