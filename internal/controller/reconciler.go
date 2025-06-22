package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// handleResourceChange processes the resource change and triggers rollouts if needed
func (r *AutoRolloutReconciler) handleResourceChange(ctx context.Context, obj AutoRolloutResource) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !hasAutoRolloutLabel(obj) {
		return ctrl.Result{}, nil
	}

	resourceType := "Unknown"
	switch obj.(type) {
	case *corev1.ConfigMap:
		resourceType = "ConfigMap"
	}

	log.Info("Resource with autorollout label has been changed",
		"type", resourceType,
		"name", obj.GetName(),
		"namespace", obj.GetNamespace(),
		"resourceVersion", obj.GetResourceVersion(),
	)

	// TODO(aditya): Find deployments using this resource
	// TODO(aditya): Trigger rollout for each deployment
	// TODO(aditya): Handle rollout errors and retries

	return ctrl.Result{}, nil
}
