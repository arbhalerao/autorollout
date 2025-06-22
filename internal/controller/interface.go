package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

// AutoRolloutResource defines the interface for resources that can trigger autorollouts
type AutoRolloutResource interface {
	metav1.Object
	runtime.Object
}

var _ AutoRolloutResource = &corev1.ConfigMap{}

// shouldProcessUpdate returns true if resource data has changed and has autorollout label
func shouldProcessUpdate(e event.UpdateEvent) bool {
	if _, isOldCM := e.ObjectOld.(*corev1.ConfigMap); isOldCM {
		if _, isNewCM := e.ObjectNew.(*corev1.ConfigMap); isNewCM {
			return shouldProcessConfigMapUpdate(e)
		}
	}

	return false
}

// hasAutoRolloutLabel checks if resource has autorollout.io=true label
func hasAutoRolloutLabel(obj AutoRolloutResource) bool {
	labels := obj.GetLabels()
	if labels == nil {
		return false
	}
	value, exists := labels["autorollout.io"]

	return exists && value == "true"
}
