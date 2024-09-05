/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "nfd.k8s-sigs.io", Version: "v1alpha1"}

	// SchemeBuilder is the scheme builder for this API.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme is a function to register this API group and version to a scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource name and returns a Group qualified GroupResource.
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&NodeFeature{},
		&NodeFeatureRule{},
		&NodeFeatureGroup{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
