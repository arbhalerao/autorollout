package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/arbhalerao/autorollout/internal/resource"
	"github.com/arbhalerao/autorollout/internal/rollout"
)

// AutoRolloutReconciler watches ConfigMaps and Secrets and triggers rollouts.
type AutoRolloutReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	resourceWatcher *resource.Watcher
	rolloutManager  *rollout.Manager
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;patch;update

func (r *AutoRolloutReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	obj, err := r.resourceWatcher.GetResource(ctx, req.NamespacedName)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if obj == nil {
		log.Info("Resource not found, might have been deleted", "namespacedName", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	return r.rolloutManager.HandleResourceChange(ctx, obj)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoRolloutReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.resourceWatcher = resource.NewWatcher(r.Client)
	r.rolloutManager = rollout.NewManager(r.Client)

	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Watches(&corev1.Secret{}, &handler.EnqueueRequestForObject{}).
		WithEventFilter(NewAutoRolloutPredicate(r.resourceWatcher)).
		Named("autorollout").
		Complete(r)
}
