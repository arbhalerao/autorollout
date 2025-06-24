package resource

import (
	"bytes"
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/arbhalerao/autorollout/internal/labels"
)

type Watcher struct {
	client.Client
}

func NewWatcher(client client.Client) *Watcher {
	return &Watcher{Client: client}
}

func (w *Watcher) GetResource(ctx context.Context, namespacedName types.NamespacedName) (AutoRolloutResource, error) {
	// Try ConfigMap first
	var cm corev1.ConfigMap
	if err := w.Get(ctx, namespacedName, &cm); err == nil {
		return &cm, nil
	}

	// Try Secret
	var secret corev1.Secret
	if err := w.Get(ctx, namespacedName, &secret); err == nil {
		return &secret, nil
	}

	return nil, nil
}

func (w *Watcher) ShouldProcessUpdate(e event.UpdateEvent) bool {
	if _, isOldCM := e.ObjectOld.(*corev1.ConfigMap); isOldCM {
		if _, isNewCM := e.ObjectNew.(*corev1.ConfigMap); isNewCM {
			return w.shouldProcessConfigMapUpdate(e)
		}
	}

	if _, isOldSecret := e.ObjectOld.(*corev1.Secret); isOldSecret {
		if _, isNewSecret := e.ObjectNew.(*corev1.Secret); isNewSecret {
			return w.shouldProcessSecretUpdate(e)
		}
	}

	return false
}

func (w *Watcher) shouldProcessConfigMapUpdate(e event.UpdateEvent) bool {
	oldCM, newCM := e.ObjectOld.(*corev1.ConfigMap), e.ObjectNew.(*corev1.ConfigMap)
	if oldCM == nil || newCM == nil {
		return false
	}

	return labels.HasAutoRolloutLabel(newCM) && w.configMapDataChanged(oldCM, newCM)
}

func (w *Watcher) shouldProcessSecretUpdate(e event.UpdateEvent) bool {
	oldSecret, newSecret := e.ObjectOld.(*corev1.Secret), e.ObjectNew.(*corev1.Secret)
	if oldSecret == nil || newSecret == nil {
		return false
	}

	return labels.HasAutoRolloutLabel(newSecret) && w.secretDataChanged(oldSecret, newSecret)
}

func (w *Watcher) configMapDataChanged(oldCM, newCM *corev1.ConfigMap) bool {
	if len(oldCM.Data) != len(newCM.Data) {
		return true
	}

	for key, newValue := range newCM.Data {
		if oldValue, exists := oldCM.Data[key]; !exists || oldValue != newValue {
			return true
		}
	}

	return false
}

func (w *Watcher) secretDataChanged(oldSecret, newSecret *corev1.Secret) bool {
	if len(oldSecret.Data) != len(newSecret.Data) {
		return true
	}

	for key, newValue := range newSecret.Data {
		if oldValue, exists := oldSecret.Data[key]; !exists || !bytes.Equal(oldValue, newValue) {
			return true
		}
	}

	return false
}
