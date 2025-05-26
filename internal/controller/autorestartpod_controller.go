/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	stablev1 "github.com/crazyfrankie/autorestart-operator/api/v1"
)

// AutoRestartPodReconciler reconciles a AutoRestartPod object
type AutoRestartPodReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=stable.crazyfrank.com,resources=autorestartpods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.crazyfrank.com,resources=autorestartpods/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stable.crazyfrank.com,resources=autorestartpods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state,
// the AutoRestartPod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *AutoRestartPodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	obj := &stablev1.AutoRestartPod{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, err
	}

	schedule, err := parseCronSchedule(obj.Spec.Schedule)
	if err != nil {
		log.Error(err, "Failed to parse cron schedule", "schedule", obj.Spec.Schedule)
		return ctrl.Result{}, err
	}

	var now time.Time
	if obj.Spec.TimeZone != "" {
		loc, err := time.LoadLocation(obj.Spec.TimeZone)
		if err != nil {
			log.Error(err, "Failed to parse timezone", "timezone", obj.Spec.TimeZone)
			return ctrl.Result{}, err
		}
		now = time.Now().In(loc)
	} else {
		now = time.Now()
	}

	nextRun := schedule.Next(now)
	if nextRun.After(now) {
		return ctrl.Result{RequeueAfter: nextRun.Sub(now)}, nil
	}

	podList := &corev1.PodList{}
	selector, _ := metav1.LabelSelectorAsSelector(&obj.Spec.Selector)
	if err = r.List(ctx, podList, client.InNamespace(req.Namespace),
		client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return ctrl.Result{}, err
	}

	for _, pod := range podList.Items {
		if err := r.Delete(ctx, &pod); err != nil {
			log.Error(err, "Failed to delete pod", "pod ", pod.Name)
		} else {
			log.Info("Restart pod", "pod ", pod.Name)
		}
	}

	obj.Status.LastRestartTime = &metav1.Time{Time: now}
	if err := r.Status().Update(ctx, obj); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Until(nextRun)}, nil
}

// parseCronSchedule parses cron expressions in various formats,
// including standard cron and cron with seconds.
func parseCronSchedule(schedule string) (cron.Schedule, error) {
	standardParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	if sch, err := standardParser.Parse(schedule); err == nil {
		return sch, nil
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	return parser.Parse(schedule)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoRestartPodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.AutoRestartPod{}).
		Named("autorestartpod").
		Complete(r)
}
