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
	RecoveryNode string
	KRef         *v1.Pod
}

type ConfigMapEntry struct {
	Id            types.UID
	Mode          string
	LabelSelector *v12.LabelSelector
	Interval      int
	LivenessProbe Probe
	Pods          []*PsPod
}

type Probe struct {
	Interval int
	Port     int
	Path     string
}

var ConfigMap = map[types.UID]*ConfigMapEntry{}

func UpdatePod(psPod *PsPod, pod *v1.Pod) {
	restarts := int32(0)
	if len(pod.Status.ContainerStatuses) >= 1 {
		restarts = pod.Status.ContainerStatuses[0].RestartCount
	}

	kRef := *pod
	psPod.Uid = string(pod.UID)
	psPod.Name = pod.Name
	psPod.Ip = pod.Status.PodIP
	psPod.HostIp = pod.Status.HostIP
	psPod.Restarts = restarts
	psPod.Deleted = pod.DeletionTimestamp != nil
	psPod.RecoveryNode = DetermineRecoveryNode(pod.Status.HostIP)
	psPod.KRef = &kRef
}
