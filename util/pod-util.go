package util

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

func GetPodsBySelector(client client.Client, ctx context.Context, labelSelector *metav1.LabelSelector) ([]v1.Pod, error) {
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, err
	}

	allPods := &v1.PodList{}
	err = client.List(ctx, allPods)
	if err != nil {
		return nil, err
	}

	var pods []v1.Pod
	for _, pod := range allPods.Items {
		if selector.Matches(labels.Set(pod.Labels)) {
			pods = append(pods, pod)
		}
	}
	return pods, nil
}

func DeletePod(c client.Client, ctx context.Context, pod *v1.Pod) {
	_ = log.FromContext(ctx)
	err := c.Delete(ctx, pod, client.GracePeriodSeconds(0))
	if err != nil {
		log.Log.Error(err, "Error while deleting Pod")
		return
	}
}

func CreatePod(client client.Client, ctx context.Context, labelSelector *metav1.LabelSelector, psPod *PsPod) string {
	_ = log.FromContext(ctx)
	replicaSets := &appsv1.ReplicaSetList{}
	err := client.List(ctx, replicaSets)
	if err != nil {
		log.Log.Error(err, "error while getting replicaSets")
		return ""
	}

	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		log.Log.Error(err, "error while transforming selector to labelSelector")
		return ""
	}
	var ourReplicaSet appsv1.ReplicaSet
	for _, replicaSet := range replicaSets.Items {
		if *replicaSet.Spec.Replicas == 0 {
			continue
		}
		if selector.Matches(labels.Set(replicaSet.Labels)) {
			ourReplicaSet = replicaSet
		}
	}

	podTemplateHash := ourReplicaSet.Spec.Template.Labels["pod-template-hash"]
	delete(ourReplicaSet.Spec.Template.Labels, "pod-template-hash")

	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    ourReplicaSet.Name + "-",
			Namespace:       ourReplicaSet.Namespace,
			Generation:      0,
			Labels:          ourReplicaSet.Spec.Template.Labels,
			OwnerReferences: nil,
		},
		Spec:   ourReplicaSet.Spec.Template.Spec,
		Status: v1.PodStatus{},
	}
	trueVar := true
	pod.Spec.ShareProcessNamespace = &trueVar
	pod.Spec.Containers[0].Image = "localhost/" + psPod.Name + "-checkpoint:latest"
	pod.Spec.Containers[0].ImagePullPolicy = "Never"
	if pod.Spec.NodeSelector == nil {
		pod.Spec.NodeSelector = make(map[string]string)
	}
	pod.Spec.NodeSelector["hitachienergy.com/internal-ip"] = psPod.RecoveryNode
	log.Log.Info("scheduling pod on node", "node", psPod.RecoveryNode)
	err = client.Create(ctx, pod)
	if err != nil {
		log.Log.Error(err, "error while creating pod")
		return ""
	}

	log.Log.Info("created Pod", "pod", pod.Name)

	RestorePods[pod.Name] = &RestorePod{
		TemplateHash: podTemplateHash,
	}

	return pod.Name
}

func DetermineRecoveryNode(currentHostIp string) string {
	var node string
	for len(KubeletAddress) <= 1 {
		log.Log.Info("KubeletAddress not initialized yet...sleeping")
		time.Sleep(time.Second)
	}
	for ip := range KubeletAddress {
		if ip == currentHostIp || ip == "192.168.56.11" {
			continue
		}
		node = ip
		break
	}

	if node == "" {
		log.Log.Info("couldn't find suitable recovery Node", "hostIp", currentHostIp)
	}

	return node
}
