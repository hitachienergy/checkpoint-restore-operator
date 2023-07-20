package controllers

import (
	"context"
	"errors"
	"hitachienergy.com/cr-operator/handlers"
	"hitachienergy.com/cr-operator/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
)

// ConfigReconciler reconciles a PullState object
type ConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, deployment)
	if err != nil {
		log.Log.Info("Error while getting deployment")
		return ctrl.Result{}, err
	}

	if deployment.GetAnnotations()["hitachienergy.com/cr_mode"] == "" {
		return ctrl.Result{}, nil
	}

	log.Log.Info("C/R requested for", "deployment", deployment.Name, "id", deployment.UID)

	partialConfig, err := parseAnnotations(deployment.GetAnnotations())
	if err != nil {
		return ctrl.Result{}, err
	}

	if _, ok := util.ConfigMap[deployment.UID]; !ok {
		util.ConfigMap[deployment.UID] = &partialConfig
	}

	probe := deployment.Spec.Template.Spec.Containers[0].LivenessProbe
	util.ConfigMap[deployment.UID].LivenessProbe = util.Probe{
		Interval: int(probe.PeriodSeconds),
		Port:     int(probe.HTTPGet.Port.IntVal),
		Path:     probe.HTTPGet.Path,
	}

	util.ConfigMap[deployment.UID].Id = deployment.UID
	util.ConfigMap[deployment.UID].LabelSelector = deployment.Spec.Selector
	util.ConfigMap[deployment.UID].Mode = partialConfig.Mode

	psPods, err := getInitialPods(r.Client, ctx, deployment.Spec.Selector)
	if err != nil {
		return ctrl.Result{}, err
	}
	util.ConfigMap[deployment.UID].Pods = psPods

	handlers.NotifyCheckpointHandler(deployment.UID)
	handlers.NotifyLivenessHandler(r.Client, deployment.UID)

	return ctrl.Result{}, nil
}

var knownPsModes = []string{"cr"}

func parseAnnotations(annotations map[string]string) (util.ConfigMapEntry, error) {
	var knownPsMode = false
	for _, mode := range knownPsModes {
		if annotations["hitachienergy.com/cr_mode"] == mode {
			knownPsMode = true
			break
		}
	}
	if !knownPsMode {
		log.Log.Info("unknown cr_mode: ", "Mode", annotations["hitachienergy.com/cr_mode"])
		return util.ConfigMapEntry{}, errors.New("unknown cr_mode")
	}

	interval, err := strconv.Atoi(annotations["hitachienergy.com/cr_interval"])
	if err != nil {
		return util.ConfigMapEntry{}, err
	}
	return util.ConfigMapEntry{
		Interval: interval,
		Mode:     "http",
	}, nil
}

func getInitialPods(client client.Client, ctx context.Context, selector *v1.LabelSelector) ([]*util.PsPod, error) {
	pods, err := util.GetPodsBySelector(client, ctx, selector)
	if err != nil {
		return nil, err
	}
	psPods := make([]*util.PsPod, len(pods))
	for i, pod := range pods {
		psPod := &util.PsPod{}
		util.UpdatePod(psPod, &pod)
		psPods[i] = psPod
	}
	return psPods, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
