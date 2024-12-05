package ws

import (
	"fmt"
	"reflect"

	"github.com/LabGroupware/go-measure-tui/internal/job"
	"github.com/LabGroupware/go-measure-tui/internal/utils"
)

type EventHandler[T any] interface {
	Handle(ws *WebSocket, data T, raw []byte)
}

type GenericEventHandler interface {
	Handle(ws *WebSocket, data any, raw []byte)
}

type EventHandlerAdapter[T any] struct {
	Handler EventHandler[T]
}

func (a EventHandlerAdapter[T]) Handle(ws *WebSocket, data any, raw []byte) {
	d, ok := data.(T)
	if !ok {
		fmt.Println("Invalid data type")
		return
	}
	a.Handler.Handle(ws, d, raw)
}

type DetailEventResponseMessageHandler struct {
	SpecificEventHandleFuncMap map[job.JobEventType]GenericEventHandler
}

func NewDetailEventResponseMessageHandler() *DetailEventResponseMessageHandler {
	return &DetailEventResponseMessageHandler{
		SpecificEventHandleFuncMap: make(map[job.JobEventType]GenericEventHandler),
	}
}

func (h *DetailEventResponseMessageHandler) RegisterHandleFunc(eventType job.JobEventType, handler GenericEventHandler) {
	h.SpecificEventHandleFuncMap[eventType] = handler
}

func (h *DetailEventResponseMessageHandler) HandleMessage(ws *WebSocket, msg *EventResponseMessage, raw []byte) error {
	switch msg.EventType {
	case EventTypesJobBegan, EventTypesJobProcessed, EventTypesJobSuccess, EventTypesJobFailed:
		var rawData EventResponseMessageWithData[job.JobRawData]
		if err := utils.UnmarshalJSON(raw, &rawData); err != nil {
			return fmt.Errorf("failed to unmarshal job event data: %v", err)
		}
		jobTypeMapper := NewJobDataTypeMapper()
		if t, ok := jobTypeMapper[rawData.Data.JobEventType]; ok {
			data := reflect.New(t).Interface()
			if err := utils.UnmarshalJSON(raw, &data); err != nil {
				return fmt.Errorf("failed to unmarshal job began data: %v", err)
			}
			if v, ok := h.SpecificEventHandleFuncMap[rawData.Data.JobEventType]; ok {
				v.Handle(ws, data, raw)
			}
		}
	}

	return nil
}
