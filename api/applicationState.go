package api

import (
	"log"
	"sync"
)

type ApplicationState struct {
	totalClicks      int
	totalClicksMutex sync.Mutex
}

func NewApplicationState() *ApplicationState {
	return &ApplicationState{
		totalClicks: 0,
	}
}

func (appState *ApplicationState) AddClick() {
	go func() {
		appState.totalClicksMutex.Lock()
		appState.totalClicks += 1
		appState.totalClicksMutex.Unlock()
		log.Printf("%v", appState.totalClicks)
	}()
}

func (appState *ApplicationState) GetClicks() int {
	return appState.totalClicks
}
