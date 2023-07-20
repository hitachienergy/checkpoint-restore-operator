/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"
	"hitachienergy.com/pull-state-operator/handlers"
	"hitachienergy.com/pull-state-operator/util"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
)

// PullStateReconciler reconciles a PullState object
type PullStateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PullState object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *PullStateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, deployment)

	if err != nil {
		log.Log.Info("Error while getting deployment")
		return ctrl.Result{}, err
	}

	if deployment.GetAnnotations()["hitachienergy.com/ps_mode"] == "" {
		return ctrl.Result{}, nil
	}

	log.Log.Info("Pull State requested for", "deployment", deployment.Name)

	partialConfig, err := parseAnnotations(deployment.GetAnnotations())
	if err != nil {
		return ctrl.Result{}, err
	}

	if _, ok := util.ConfigMap[deployment.UID]; !ok {
		util.ConfigMap[deployment.UID] = &partialConfig
	}

	pods, err := r.handlePods(ctx, deployment)
	if err != nil {
		return ctrl.Result{}, err
	}

	probe := deployment.Spec.Template.Spec.Containers[0].LivenessProbe
	util.ConfigMap[deployment.UID].LivenessProbe = util.Probe{
		Interval: int(probe.PeriodSeconds),
		Port:     int(probe.HTTPGet.Port.IntVal),
		Path:     probe.HTTPGet.Path,
	}

	util.ConfigMap[deployment.UID].Id = deployment.UID
	util.ConfigMap[deployment.UID].StateProbe = partialConfig.StateProbe
	util.ConfigMap[deployment.UID].LabelSelector = deployment.Spec.Selector
	util.ConfigMap[deployment.UID].Mode = partialConfig.Mode
	util.ConfigMap[deployment.UID].Pods = pods

	handlers.NotifyStateHandler(deployment.UID)
	if partialConfig.Mode == "http-sf" {
		log.Log.Info("starting liveness handler")
		handlers.NotifyLivenessHandler(r.Client, deployment.UID)
	}

	return ctrl.Result{}, nil
}

func (r *PullStateReconciler) handlePods(ctx context.Context, deployment *appsv1.Deployment) ([]*util.PsPod, error) {
	kPods, err := util.GetPodsBySelector(r.Client, ctx, deployment.Spec.Selector)
	if err != nil {
		return nil, err
	}

	var pods []*util.PsPod
	for _, pod := range kPods {
		if pod.Status.PodIP == "" {
			// Do not add Pods, that are not yet scheduled
			continue
		}

		var lastPod *util.PsPod
		for _, psPod := range util.ConfigMap[deployment.UID].Pods {
			if psPod.Uid == string(pod.UID) {
				lastPod = psPod
			}
		}
		currentPod := util.FromPod(&pod)
		if lastPod != nil && currentPod.Restarts-lastPod.Restarts == 1 {
			log.Log.Info("Container in Pod restarted", "pod", currentPod.Name)
			util.Restores <- util.Restore{
				FromPid: currentPod.Uid,
				HostIp:  currentPod.HostIp,
				Ip:      currentPod.Ip,
				Mode:    util.ConfigMap[deployment.UID].Mode,
				Path:    util.ConfigMap[deployment.UID].StateProbe.Path,
				Port:    int32(util.ConfigMap[deployment.UID].StateProbe.Port),
			}
		}

		if lastPod != nil && !lastPod.Deleted && currentPod.Deleted {
			log.Log.Info("Pod was scheduled for deletion in this round", "pod", currentPod.Name)
			util.GetDeletedPodsChannelFor(deployment.UID) <- currentPod
		}

		if lastPod == nil {
			log.Log.Info("New Pod arrived", "pod", currentPod.Name)
			// if there is a new pod and a deleted pod in the same update, we want to first schedule
			// the deleted pod and then schedule the new pod.
			defer func() { util.GetNewPodsChannelFor(deployment.UID) <- currentPod }()
		}

		pods = append(pods, &currentPod)

		log.Log.Info("POD INFO:",
			"pod", pod.Name,
			"ip", pod.Status.PodIP,
			"node", pod.Status.NominatedNodeName,
			"hostIp", pod.Status.HostIP,
		)
	}
	return pods, nil
}

var knownPsModes = []string{"http", "http-sf"}

func parseAnnotations(annotations map[string]string) (util.ConfigMapEntry, error) {
	var knownPsMode = false
	for _, mode := range knownPsModes {
		if annotations["hitachienergy.com/ps_mode"] == mode {
			knownPsMode = true
			break
		}
	}
	if !knownPsMode {
		log.Log.Info("unknown ps_mode: ", "Mode", annotations["hitachienergy.com/ps_mode"])
		return util.ConfigMapEntry{}, errors.New("unknown ps_mode")
	}

	interval, err := strconv.Atoi(annotations["hitachienergy.com/ps_interval"])
	if err != nil {
		return util.ConfigMapEntry{}, err
	}
	port, err := strconv.Atoi(annotations["hitachienergy.com/ps_port"])
	if err != nil {
		return util.ConfigMapEntry{}, err
	}
	path := annotations["hitachienergy.com/ps_path"]
	return util.ConfigMapEntry{
		StateProbe: util.Probe{
			Interval: interval,
			Port:     port,
			Path:     path,
		},
		Mode: annotations["hitachienergy.com/ps_mode"],
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PullStateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
