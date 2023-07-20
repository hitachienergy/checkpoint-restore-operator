package util

// KubeletPorts map of node ip to kubelet port
var KubeletPorts = make(map[string]int32)

// KubeletAddress map of node ip to kubelet address (ip or hostname)
var KubeletAddress = make(map[string]string)
