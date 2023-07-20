package controllers

import (
	"context"
	"hitachienergy.com/pull-state-operator/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type NodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;patch

func (r *NodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	node := &v1.Node{}
	err := r.Get(ctx, req.NamespacedName, node)
	if err != nil {
		log.Log.Error(err, "error while getting node", "node", req.Name)
		return ctrl.Result{}, nil
	}

	var ip string
	for _, address := range node.Status.Addresses {
		if address.Type == "InternalIP" {
			ip = address.Address
		}
	}
	util.Nodes[ip] = struct{}{}
	// TODO: detect missing nodes
	patch := client.MergeFrom(node.DeepCopy())
	node.Labels["hitachienergy.com/internal-ip"] = ip
	err = r.Patch(ctx, node, patch)
	if err != nil {
		log.Log.Error(err, "unable to set node label for ip", "node", node.Name)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Node{}).
		Complete(r)
}
