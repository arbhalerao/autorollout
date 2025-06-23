package controller

import (
	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// TODO(aditya): Trigger rollout for each deployment
// TODO(aditya): Handle rollout errors and retries
func (r *AutoRolloutReconciler) handleResourceChange(ctx context.Context, obj AutoRolloutResource) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !r.hasAutoRolloutLabel(obj) {
		return ctrl.Result{}, nil
	}

	resourceType := "Unknown"
	var err error

	switch resource := obj.(type) {
	case *corev1.ConfigMap:
		resourceType = "ConfigMap"
		err = r.handleConfigMapChange(ctx, resource)
		if err != nil {
			log.Error(err, "Failed to handle ConfigMap change")
			return ctrl.Result{RequeueAfter: time.Minute * 5}, err
		}
	case *corev1.Secret:
		resourceType = "Secret"
		err = r.handleSecretChange(ctx, resource)
		if err != nil {
			log.Error(err, "Failed to handle Secret change")
			return ctrl.Result{RequeueAfter: time.Minute * 5}, err
		}
	}

	log.Info("Resource with autorollout label has been changed",
		"type", resourceType,
		"name", obj.GetName(),
		"namespace", obj.GetNamespace(),
		"resourceVersion", obj.GetResourceVersion(),
	)

	return ctrl.Result{}, nil
}

func (r *AutoRolloutReconciler) handleConfigMapChange(ctx context.Context, cm *corev1.ConfigMap) error {
	log := logf.FromContext(ctx)

	deployments, err := r.findDeploymentsUsingConfigMap(ctx, cm)
	if err != nil {
		log.Error(err, "Failed to find deployments using ConfigMap")
		return err
	}

	if len(deployments) == 0 {
		log.Info("No deployments found using this ConfigMap", "configmap", cm.Name)
		return nil
	}

	for _, deployment := range deployments {
		log.Info("Deployment using ConfigMap", "deployment", deployment.Name)
	}

	return nil
}

func (r *AutoRolloutReconciler) handleSecretChange(ctx context.Context, secret *corev1.Secret) error {
	log := logf.FromContext(ctx)

	deployments, err := r.findDeploymentsUsingSecret(ctx, secret)
	if err != nil {
		log.Error(err, "Failed to find deployments using Secret")
		return err
	}

	if len(deployments) == 0 {
		log.Info("No deployments found using this Secret", "secret", secret.Name)
		return nil
	}

	for _, deployment := range deployments {
		log.Info("Deployment using Secret", "deployment", deployment.Name)
	}

	return nil
}
