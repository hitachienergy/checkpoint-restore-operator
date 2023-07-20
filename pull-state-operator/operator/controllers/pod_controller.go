package controllers

import (
	"context"
	"hitachienergy.com/pull-state-operator/util"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			log.Log.Info("not found", "err", err)
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "error while getting pod", "pod", req.Name)
		return ctrl.Result{}, nil
	}
	for _, entry := range util.ConfigMap {
		selector, err := v12.LabelSelectorAsSelector(entry.LabelSelector)
		if err != nil {
			log.Log.Error(err, "error while converting labelSelector")
			continue
		}

		if selector.Matches(labels.Set(pod.Labels)) {
			log.Log.Info("Update for Pod", "pod", pod.Name, "host", pod.Status.HostIP)
			if podTemplateHash, ok := util.RestorePods[pod.Name]; ok {
				if pod.Status.PodIP != "" {
					log.Log.Info("Ip is now known!", "ip", pod.Status.PodIP)
					util.GetNewPodsChannelFor(entry.Id) <- util.FromPod(pod)
					patch := client.MergeFrom(pod.DeepCopy())
					pod.Labels["pod-template-hash"] = podTemplateHash
					err := r.Patch(ctx, pod, patch)
					if err != nil {
						log.Log.Error(err, "updating pod template hash failed!")
						return ctrl.Result{}, nil
					}
					delete(util.RestorePods, pod.Name)
				}
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		Complete(r)
}
