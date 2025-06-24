package resource

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func (w *Watcher) FindDeploymentsUsingConfigMap(ctx context.Context, cm *corev1.ConfigMap) ([]appsv1.Deployment, error) {
	log := logf.FromContext(ctx)

	deploymentList := &appsv1.DeploymentList{}
	err := w.List(ctx, deploymentList, client.InNamespace(cm.GetNamespace()))
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	var usingDeployments []appsv1.Deployment
	cmName := cm.GetName()

	for _, dep := range deploymentList.Items {
		if w.deploymentUsesConfigMap(&dep, cmName) {
			usingDeployments = append(usingDeployments, dep)
			log.Info("Found deployment using ConfigMap",
				"deployment", dep.Name,
				"configmap", cm.Name,
			)
		}
	}

	return usingDeployments, nil
}

func (w *Watcher) FindDeploymentsUsingSecret(ctx context.Context, secret *corev1.Secret) ([]appsv1.Deployment, error) {
	log := logf.FromContext(ctx)

	deploymentList := &appsv1.DeploymentList{}
	err := w.List(ctx, deploymentList, client.InNamespace(secret.GetNamespace()))
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	var usingDeployments []appsv1.Deployment
	secretName := secret.GetName()

	for _, dep := range deploymentList.Items {
		if w.deploymentUsesSecret(&dep, secretName) {
			usingDeployments = append(usingDeployments, dep)
			log.Info("Found deployment using Secret",
				"deployment", dep.Name,
				"secret", secret.Name,
			)
		}
	}

	return usingDeployments, nil
}

func (w *Watcher) FindAffectedDeployments(ctx context.Context, resource AutoRolloutResource) ([]appsv1.Deployment, error) {
	switch res := resource.(type) {
	case *corev1.ConfigMap:
		return w.FindDeploymentsUsingConfigMap(ctx, res)
	case *corev1.Secret:
		return w.FindDeploymentsUsingSecret(ctx, res)
	default:
		return nil, fmt.Errorf("unsupported resource type: %T", resource)
	}
}

func (w *Watcher) deploymentUsesConfigMap(dep *appsv1.Deployment, cmName string) bool {
	podSpec := dep.Spec.Template.Spec

	for _, vol := range podSpec.Volumes {
		if vol.ConfigMap != nil && vol.ConfigMap.Name == cmName {
			return true
		}
	}

	for _, container := range podSpec.Containers {
		if w.containerUsesConfigMap(&container, cmName) {
			return true
		}
	}

	for _, container := range podSpec.InitContainers {
		if w.containerUsesConfigMap(&container, cmName) {
			return true
		}
	}

	return false
}

func (w *Watcher) deploymentUsesSecret(dep *appsv1.Deployment, secretName string) bool {
	podSpec := dep.Spec.Template.Spec

	for _, vol := range podSpec.Volumes {
		if vol.Secret != nil && vol.Secret.SecretName == secretName {
			return true
		}
	}

	for _, imagePullSecret := range podSpec.ImagePullSecrets {
		if imagePullSecret.Name == secretName {
			return true
		}
	}

	for _, container := range podSpec.Containers {
		if w.containerUsesSecret(&container, secretName) {
			return true
		}
	}

	for _, container := range podSpec.InitContainers {
		if w.containerUsesSecret(&container, secretName) {
			return true
		}
	}

	return false
}

func (w *Watcher) containerUsesConfigMap(container *corev1.Container, cmName string) bool {
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

	return false
}

func (w *Watcher) containerUsesSecret(container *corev1.Container, secretName string) bool {
	for _, env := range container.Env {
		if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
			if env.ValueFrom.SecretKeyRef.Name == secretName {
				return true
			}
		}
	}

	for _, envFrom := range container.EnvFrom {
		if envFrom.SecretRef != nil && envFrom.SecretRef.Name == secretName {
			return true
		}
	}

	return false
}
