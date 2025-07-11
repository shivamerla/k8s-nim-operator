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
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1alpha1 "github.com/NVIDIA/k8s-nim-operator/api/apps/v1alpha1"
)

var _ = Describe("DistributedNIMService Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "distributed-service"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		distributednimservice := &appsv1alpha1.DistributedNIMService{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind DistributedNIMService")
			err := k8sClient.Get(ctx, typeNamespacedName, distributednimservice)
			if err != nil && errors.IsNotFound(err) {
				resource := &appsv1alpha1.DistributedNIMService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: appsv1alpha1.DistributedNIMServiceSpec{
						Image: appsv1alpha1.Image{
							Repository: "nvcr.io/nvidia/ai-dynamo/vllm-runtime",
							Tag:        "0.3.1",
						},
						AuthSecret: "ngc-secret",
						Storage: appsv1alpha1.NIMServiceStorage{
							NIMCache: appsv1alpha1.NIMCacheVolSpec{Name: "test-pvc"},
						},
						Expose: appsv1alpha1.Expose{
							Service: appsv1alpha1.Service{Type: "ClusterIP", Port: pointer.Int32(8000)},
						},
						Config: &appsv1alpha1.DynamoConfig{
							Name:      "disagg-config",
							MountPath: "/opt/config",
						},
						Workers: &appsv1alpha1.WorkerSpec{
							Prefill: &appsv1alpha1.PrefillSpec{
								Replicas: 1,
							},
							Decode: &appsv1alpha1.DecoderSpec{
								Replicas: 1,
							},
						},
						Router: &appsv1alpha1.RouterSpec{
							Replicas: 1,
						},
						Planner: &appsv1alpha1.PlannerSpec{
							Replicas: 1,
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &appsv1alpha1.DistributedNIMService{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance DistributedNIMService")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &DistributedNIMServiceReconciler{
				Client: k8sClient,
				scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
