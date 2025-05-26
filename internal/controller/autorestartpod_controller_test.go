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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	stablev1 "github.com/crazyfrankie/autorestart-operator/api/v1"
)

var _ = Describe("AutoRestartPod Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "kubectl-system",
		}
		autorestartpod := &stablev1.AutoRestartPod{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind AutoRestartPod")
			err := k8sClient.Get(ctx, typeNamespacedName, autorestartpod)
			if err != nil && errors.IsNotFound(err) {
				resource := &stablev1.AutoRestartPod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "kubectl-system",
					},
					Spec: stablev1.AutoRestartPodSpec{
						Schedule: "* 0/5 * * * ? ",
						Selector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx",
							},
						},
						TimeZone: "Asia/Shanghai",
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &stablev1.AutoRestartPod{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance AutoRestartPod")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &AutoRestartPodReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify the resource was updated correctly
			updatedResource := &stablev1.AutoRestartPod{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedResource)
			Expect(err).NotTo(HaveOccurred())

			// Check if LastRestartTime was set, which indicates successful reconciliation
			Expect(updatedResource.Status.LastRestartTime).NotTo(BeNil(), "LastRestartTime should be set after reconciliation")

			// Verify the requeue time is set for the next scheduled run
			// Note: This test is simplified, in a real test you might mock time.Now() to control timing
		})
	})
})
