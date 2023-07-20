package util

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Restore struct {
	FromPid string
	HostIp  string
	Ip      string
	Mode    string
	Path    string
	Port    int32
}

var RestorePods = make(map[string]string)

var DeletedPods = map[types.UID]chan PsPod{}
var NewPods = map[types.UID]chan PsPod{}

var Restores = make(chan Restore)

func GetDeletedPodsChannelFor(deployment types.UID) chan PsPod {
	if DeletedPods[deployment] == nil {
		makeChannels(deployment)
	}
	return DeletedPods[deployment]
}

func GetNewPodsChannelFor(deployment types.UID) chan PsPod {
	if NewPods[deployment] == nil {
		makeChannels(deployment)
	}
	return NewPods[deployment]
}

func makeChannels(deployment types.UID) {
	DeletedPods[deployment] = make(chan PsPod, 10)
	NewPods[deployment] = make(chan PsPod, 10)
	go CombineLatest(deployment)
}

func CombineLatest(deployment types.UID) {
	var nPod PsPod
	var dPod PsPod
	for {
		select {
		case dPod = <-DeletedPods[deployment]:
			nPod = <-NewPods[deployment]
			log.Log.Info("Scheduling Restore", "from", dPod.Name, "to", nPod.Name)
			Restores <- Restore{
				FromPid: dPod.Uid,
				HostIp:  nPod.HostIp,
				Ip:      nPod.Ip,
				Mode:    ConfigMap[deployment].Mode,
				Path:    ConfigMap[deployment].StateProbe.Path,
				Port:    int32(ConfigMap[deployment].StateProbe.Port),
			}
			break
		case nPod = <-NewPods[deployment]:
			log.Log.Info("ignoring new pod", "pod", nPod.Name)
			break
		}
	}
}
