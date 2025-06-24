package labels

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AutoRolloutLabel = "autorollout.io"

func HasAutoRolloutLabel(obj metav1.Object) bool {
	labels := obj.GetLabels()
	if labels == nil {
		return false
	}
	value, exists := labels[AutoRolloutLabel]

	return exists && value == "true"
}
