package storage

import (
	"hitachienergy.com/pull-state-operator/helper/generated"
)

var stateStorage = make(map[string]*generated.State)

func StoreState(state *generated.State) {
	stateStorage[state.Id] = state
}

func GetState(id string) *generated.State {
	return stateStorage[id]
}

func Delete(pod *generated.PodId) {
	delete(stateStorage, pod.Id)
}
