package controller

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// shouldProcessConfigMapUpdate returns true if ConfigMap data has changed and has autorollout label
func (r *AutoRolloutReconciler) shouldProcessConfigMapUpdate(e event.UpdateEvent) bool {
	oldCM, newCM := e.ObjectOld.(*corev1.ConfigMap), e.ObjectNew.(*corev1.ConfigMap)
	if oldCM == nil || newCM == nil {
		return false
	}

	return r.hasAutoRolloutLabel(newCM) && r.configMapDataChanged(oldCM, newCM)
}

// configMapDataChanged compares the Data field between old and new ConfigMap
func (r *AutoRolloutReconciler) configMapDataChanged(oldCM, newCM *corev1.ConfigMap) bool {
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
