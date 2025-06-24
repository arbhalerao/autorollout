package resource

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// AutoRolloutResource defines the interface for resources that can trigger autorollouts
type AutoRolloutResource interface {
	metav1.Object
	runtime.Object
}

var _ AutoRolloutResource = &corev1.ConfigMap{}
var _ AutoRolloutResource = &corev1.Secret{}
