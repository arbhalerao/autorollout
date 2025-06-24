package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/arbhalerao/autorollout/internal/resource"
)

func NewAutoRolloutPredicate(watcher *resource.Watcher) predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return watcher.ShouldProcessUpdate(e)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
	}
}
