package controllers

import (
	"context"
	"hitachienergy.com/cr-operator/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	pod := &v1.Pod{}
	err := r.Get(ctx, req.NamespacedName, pod)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Log.Info("Pod not found. ignoring", "pod", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		log.Log.Info("Error while getting pod")
		return ctrl.Result{}, err
	}

	for _, conf := range util.ConfigMap {
		selector, err := metav1.LabelSelectorAsSelector(conf.LabelSelector)
		if err != nil {
			return ctrl.Result{}, err
		}

		if selector.Matches(labels.Set(pod.Labels)) {
			log.Log.Info("New Pod found for ", "confId", string(conf.Id), "pod", pod.Name)
			var confPod *util.PsPod
			for _, psPod := range conf.Pods {
				if string(pod.UID) == psPod.Uid {
					confPod = psPod
					break
				}
			}
			if confPod == nil {
				confPod = &util.PsPod{}
			}
			util.UpdatePod(confPod, pod)
		}
	}

	if restorePod, ok := util.RestorePods[pod.Name]; ok && pod.Status.PodIP != "" {
		log.Log.Info("Ip is now known!", "ip", pod.Status.PodIP)
		patch := client.MergeFrom(pod.DeepCopy())
		pod.Labels["pod-template-hash"] = restorePod.TemplateHash
		util.DeletePod(r.Client, ctx, restorePod.OldPod)
		err := r.Patch(ctx, pod, patch)
		if err != nil {
			log.Log.Error(err, "updating pod template hash failed!")
			return ctrl.Result{}, nil
		}
		delete(util.RestorePods, pod.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controllers with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		Complete(r)
}
