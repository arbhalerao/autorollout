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
var _ AutoRolloutResource = &corev1.Secret{}

func (r *AutoRolloutReconciler) shouldProcessUpdate(e event.UpdateEvent) bool {
	if _, isOldCM := e.ObjectOld.(*corev1.ConfigMap); isOldCM {
		if _, isNewCM := e.ObjectNew.(*corev1.ConfigMap); isNewCM {
			return r.shouldProcessConfigMapUpdate(e)
		}
	}

	if _, isOldSecret := e.ObjectOld.(*corev1.Secret); isOldSecret {
		if _, isNewSecret := e.ObjectNew.(*corev1.Secret); isNewSecret {
			return r.shouldProcessSecretUpdate(e)
		}
	}

	return false
}

func (r *AutoRolloutReconciler) hasAutoRolloutLabel(obj AutoRolloutResource) bool {
	labels := obj.GetLabels()
	if labels == nil {
		return false
	}
	value, exists := labels["autorollout.io"]

	return exists && value == "true"
}
