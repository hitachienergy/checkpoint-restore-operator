package util

import v1 "k8s.io/api/core/v1"

type RestorePod struct {
	TemplateHash string
	OldPod       *v1.Pod
}

var RestorePods = make(map[string]*RestorePod)
