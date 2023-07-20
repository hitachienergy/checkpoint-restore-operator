package util

import (
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PsPod struct {
	Uid          string
	Name         string
	Ip           string
	HostIp       string
	Restarts     int32
	Deleted      bool
	CurrentState []byte
	ContentType  string
	KRef         *v1.Pod
}

type ConfigMapEntry struct {
	Id            types.UID
	Mode          string
	LabelSelector *v12.LabelSelector
	StateProbe    Probe
	LivenessProbe Probe
	Pods          []*PsPod
}

type Probe struct {
	Interval int
	Port     int
	Path     string
}

var ConfigMap = map[types.UID]*ConfigMapEntry{}

func FromPod(pod *v1.Pod) PsPod {
	restarts := int32(0)
	if len(pod.Status.ContainerStatuses) >= 1 {
		restarts = pod.Status.ContainerStatuses[0].RestartCount
	}

	kRef := *pod
	return PsPod{
		Uid:      string(pod.UID),
		Name:     pod.Name,
		Ip:       pod.Status.PodIP,
		HostIp:   pod.Status.HostIP,
		Restarts: restarts,
		Deleted:  pod.DeletionTimestamp != nil,
		KRef:     &kRef,
	}
}
