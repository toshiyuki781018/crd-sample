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

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	stablev1 "example.com/sampleweb/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FirstCrdReconciler reconciles a FirstCrd object
type FirstCrdReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=stable.example.com,resources=firstcrds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stable.example.com,resources=firstcrds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stable.example.com,resources=firstcrds/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FirstCrd object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.0/pkg/reconcile
func (r *FirstCrdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	// your logic here
	fcrd := stablev1.FirstCrd{}
	if err := r.Get(ctx, req.NamespacedName, &fcrd); err != nil {
		log.Error(err, "unable to fetch FirstCrd")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploymentTemplate := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fcrd.Name,
			Namespace: req.Namespace,
		},
	}
	_, err := ctrl.CreateOrUpdate(ctx, r.Client, deploymentTemplate, func() error {
		ls := map[string]string{"app": "firstcrd", "firstcrd_cr": fcrd.Name}
		replicas := fcrd.Spec.Replicas

		deploymentTemplate.Spec = appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "nginx",
						Name:  "nginx",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "nginx",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "shared",
							MountPath: "/usr/share/nginx/html",
						}},
					}},
					InitContainers: []corev1.Container{{
						Image: "busybox",
						Name:  "init",
						Args:  []string{"/bin/sh", "-c", "echo" + " " + fcrd.Spec.Message + " " + ">" + " " + "/test/index.html"},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "shared",
							MountPath: "/test",
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "shared",
					}},
				},
			},
		}
		if err := ctrl.SetControllerReference(&fcrd, deploymentTemplate, r.Scheme); err != nil {
			log.Error(err, "unable to set ownerReference from FirstCrd to Deployment")
			return err
		}
		return nil
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	// FirstCrdが所有者のdeploymentを取得
	var deployment appsv1.Deployment
	var deploymentNamespacedName = client.ObjectKey{Namespace: req.Namespace, Name: fcrd.Name}
	if err := r.Get(ctx, deploymentNamespacedName, &deployment); err != nil {
		log.Error(err, "unable to fetch Deployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// FirstCrdのCurrentReplicasステータスの更新要否を確認
	updateNeeded := false
	availableReplicas := deployment.Status.AvailableReplicas
	if availableReplicas != fcrd.Status.CurrentReplicas {
		fcrd.Status.CurrentReplicas = availableReplicas
		updateNeeded = true
	}

	currentMessage := fcrd.Status.CurrentMessage
	if currentMessage != fcrd.Spec.Message {
		fcrd.Status.CurrentMessage = fcrd.Spec.Message
		updateNeeded = true
	}

	// FirstCrdのCurrentReplicasステータスを更新
	if updateNeeded {
		if err := r.Status().Update(ctx, &fcrd); err != nil {
			log.Error(err, "unable to update FirstCrd status")
			return ctrl.Result{}, err
		}
		r.Recorder.Eventf(&fcrd, corev1.EventTypeNormal, "Updated", "Update FirstCrd Status")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FirstCrdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stablev1.FirstCrd{}).
		Named("firstcrd").
		Complete(r)
}
