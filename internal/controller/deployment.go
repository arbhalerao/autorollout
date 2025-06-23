package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *AutoRolloutReconciler) findDeploymentsUsingConfigMap(ctx context.Context, cm *corev1.ConfigMap) ([]appsv1.Deployment, error) {
	log := logf.FromContext(ctx)

	deploymentList := &appsv1.DeploymentList{}
	err := r.List(ctx, deploymentList, client.InNamespace(cm.GetNamespace()))
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	var usingDeployments []appsv1.Deployment
	cmName := cm.GetName()

	for _, dep := range deploymentList.Items {
		if r.deploymentUsesConfigMap(&dep, cmName) {
			usingDeployments = append(usingDeployments, dep)
			log.Info("Found deployment using ConfigMap",
				"deployment", dep.Name,
				"configmap", cm.Name,
			)
		}
	}

	return usingDeployments, nil
}

func (r *AutoRolloutReconciler) deploymentUsesConfigMap(dep *appsv1.Deployment, cmName string) bool {
	podSpec := dep.Spec.Template.Spec

	for _, vol := range podSpec.Volumes {
		if vol.ConfigMap != nil && vol.ConfigMap.Name == cmName {
			return true
		}
	}

	for _, container := range podSpec.Containers {
		for _, env := range container.Env {
			if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil {
				if env.ValueFrom.ConfigMapKeyRef.Name == cmName {
					return true
				}
			}
		}
		for _, envFrom := range container.EnvFrom {
			if envFrom.ConfigMapRef != nil && envFrom.ConfigMapRef.Name == cmName {
				return true
			}
		}
	}

	for _, container := range podSpec.InitContainers {
		for _, env := range container.Env {
			if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil {
				if env.ValueFrom.ConfigMapKeyRef.Name == cmName {
					return true
				}
			}
		}
		for _, envFrom := range container.EnvFrom {
			if envFrom.ConfigMapRef != nil && envFrom.ConfigMapRef.Name == cmName {
				return true
			}
		}
	}

	return false
}
