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
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NamespaceReconciler", func() {
	BeforeEach(func() {
		// Start the controller inside a goroutine
		go func() {
			err := (&NamespaceReconciler{
				Client: k8sManager.GetClient(),
				Scheme: k8sManager.GetScheme(),
			}).SetupWithManager(k8sManager)
			Expect(err).NotTo(HaveOccurred())

			Expect(k8sManager.Start(ctx)).To(Succeed())
		}()
	})

	It("should add default labels to namespaces missing them", func() {
		ns := &corev1.Namespace{}
		ns.Name = "test-ns"
		Expect(k8sClient.Create(ctx, ns)).To(Succeed())

		// Defer cleanup
		DeferCleanup(func() {
			_ = k8sClient.Delete(ctx, ns)
		})

		Eventually(func() map[string]string {
			var updated corev1.Namespace
			err := k8sClient.Get(ctx, types.NamespacedName{Name: ns.Name}, &updated)
			if err != nil {
				return nil
			}
			return updated.Labels
		}, 5*time.Second, 500*time.Millisecond).Should(SatisfyAll(
			HaveKeyWithValue("team", "unknown"),
			HaveKeyWithValue("env", "dev"),
		))
	})
})
