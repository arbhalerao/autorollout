package controller

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func (r *AutoRolloutReconciler) shouldProcessSecretUpdate(e event.UpdateEvent) bool {
	oldSecret, newSecret := e.ObjectOld.(*corev1.Secret), e.ObjectNew.(*corev1.Secret)
	if oldSecret == nil || newSecret == nil {
		return false
	}

	return r.hasAutoRolloutLabel(newSecret) && r.secretDataChanged(oldSecret, newSecret)
}

func (r *AutoRolloutReconciler) secretDataChanged(oldSecret, newSecret *corev1.Secret) bool {
	if len(oldSecret.Data) != len(newSecret.Data) {
		return true
	}

	for key, newValue := range newSecret.Data {
		if oldValue, exists := oldSecret.Data[key]; !exists || !bytesEqual(oldValue, newValue) {
			return true
		}
	}

	return false
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
