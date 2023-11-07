package api

type ApplicationState struct {
	currentItemIndex int
}

func NewApplicationState() *ApplicationState {
	return &ApplicationState{}
}

func (appState *ApplicationState) AddItem() int {
	appState.currentItemIndex += 1
	return appState.currentItemIndex
}

func (appState *ApplicationState) DeleteItem() {
}

func (appState *ApplicationState) DeleteAll() {
	appState.currentItemIndex = 0
}
