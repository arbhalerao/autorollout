package rollout

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/arbhalerao/autorollout/internal/labels"
	"github.com/arbhalerao/autorollout/internal/resource"
)

type Manager struct {
	client.Client
	resourceWatcher *resource.Watcher
}

func NewManager(client client.Client) *Manager {
	return &Manager{
		Client:          client,
		resourceWatcher: resource.NewWatcher(client),
	}
}

func (m *Manager) HandleResourceChange(ctx context.Context, obj resource.AutoRolloutResource) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if !labels.HasAutoRolloutLabel(obj) {
		return ctrl.Result{}, nil
	}

	resourceType := "Unknown"
	var err error

	switch resource := obj.(type) {
	case *corev1.ConfigMap:
		resourceType = "ConfigMap"
		err = m.handleConfigMapChange(ctx, resource)
		if err != nil {
			log.Error(err, "Failed to handle ConfigMap change")
			return ctrl.Result{RequeueAfter: time.Minute * 5}, err
		}
	case *corev1.Secret:
		resourceType = "Secret"
		err = m.handleSecretChange(ctx, resource)
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

func (m *Manager) handleConfigMapChange(ctx context.Context, cm *corev1.ConfigMap) error {
	log := logf.FromContext(ctx)

	deployments, err := m.resourceWatcher.FindDeploymentsUsingConfigMap(ctx, cm)
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

	return m.rolloutDeployments(ctx, deployments)
}

func (m *Manager) handleSecretChange(ctx context.Context, secret *corev1.Secret) error {
	log := logf.FromContext(ctx)

	deployments, err := m.resourceWatcher.FindDeploymentsUsingSecret(ctx, secret)
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

	return m.rolloutDeployments(ctx, deployments)
}

func (m *Manager) rolloutDeployments(ctx context.Context, deployments []appsv1.Deployment) error {
	log := logf.FromContext(ctx)

	for _, deployment := range deployments {
		log.Info("Triggering rollout for deployment", "deployment", deployment.Name)

		updated := deployment.DeepCopy()
		if updated.Spec.Template.Annotations == nil {
			updated.Spec.Template.Annotations = map[string]string{}
		}

		updated.Spec.Template.Annotations["autorolloutTimestamp"] = time.Now().Format(time.RFC3339)

		if err := m.Update(ctx, updated); err != nil {
			log.Error(err, "Failed to trigger rollout for deployment", "deployment", deployment.Name)
			return err
		}
	}

	return nil
}
