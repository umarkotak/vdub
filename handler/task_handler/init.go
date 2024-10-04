package task_handler

import "sync"

type (
	HandlerState struct {
		mu          sync.Mutex
		RunningTask map[string]bool
	}
)

var (
	handlerState HandlerState
)

func Initialize() {
	handlerState = HandlerState{
		RunningTask: map[string]bool{},
	}
}
