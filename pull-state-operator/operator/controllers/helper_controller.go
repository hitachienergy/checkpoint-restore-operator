package controllers

import (
	"context"
	helperclient "hitachienergy.com/pull-state-operator/helper-client"
	"hitachienergy.com/pull-state-operator/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type HelperReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

func (r *HelperReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO: find a way to detect this dynamically
	if req.Namespace != "pull-state-operator-system" {
		// Skip this as we only want to watch the helper daemonSet
		return ctrl.Result{}, nil
	}
	daemonSet := &appsv1.DaemonSet{}
	err := r.Get(ctx, req.NamespacedName, daemonSet)
	if err != nil {
		log.Log.Error(err, "Error while getting daemonSet")
		return ctrl.Result{}, err
	}

	if daemonSet.Spec.Selector.MatchLabels["app"] != "ps-helper" {
		// Skip this as we only want to watch the helper daemonSet
		return ctrl.Result{}, nil
	}

	kPods, err := util.GetPodsBySelector(r.Client, ctx, daemonSet.Spec.Selector)
	if err != nil {
		log.Log.Error(err, "Error while getting pods")
		return ctrl.Result{}, err
	}

	var ips = make([]helperclient.Helper, 0, len(kPods))
	for _, pod := range kPods {
		var podReady = false
		for _, condition := range pod.Status.Conditions {
			if condition.Type == "Ready" {
				podReady = condition.Status == "True"
				break
			}
		}
		if pod.Status.Phase == v1.PodSucceeded ||
			pod.Status.Phase == v1.PodFailed ||
			pod.Status.PodIP == "" ||
			pod.Status.HostIP == "" ||
			!podReady {
			continue
		}

		helper := helperclient.Helper{
			Ip:     pod.Status.PodIP,
			HostIp: pod.Status.HostIP,
		}
		ips = append(ips, helper)

		if helperclient.HelperExists(helper) {
			continue
		}
	}
	helperclient.SetHelpers(&ips)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelperReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
		Complete(r)
}
