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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DistributedNIMServiceSpec defines the desired state of DistributedNIMService
type DistributedNIMServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Global graph parameters
	Image Image `json:"image"`
	// The name of an existing pull secret containing the NGC_API_KEY
	AuthSecret       string            `json:"authSecret"`
	Storage          NIMServiceStorage `json:"storage"`
	SchedulerName    string            `json:"schedulerName,omitempty"`
	UserID           *int64            `json:"userID,omitempty"`
	GroupID          *int64            `json:"groupID,omitempty"`
	RuntimeClassName string            `json:"runtimeClassName,omitempty"`
	Proxy            *ProxySpec        `json:"proxy,omitempty"`
	Config           *DynamoConfig     `json:"config"`

	// Frontend parameters
	Expose       Expose                `json:"expose"`
	Args         []string              `json:"args,omitempty"`
	Env          []corev1.EnvVar       `json:"env,omitempty"`
	Resources    *ResourceRequirements `json:"resources,omitempty"`
	Annotations  map[string]string     `json:"annotations,omitempty"`
	NodeSelector map[string]string     `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration   `json:"tolerations,omitempty"`
	PodAffinity  *corev1.PodAffinity   `json:"podAffinity,omitempty"`
	Metrics      Metrics               `json:"metrics,omitempty"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default:=1
	Replicas int `json:"replicas,omitempty"`

	// Other graph components
	Workers *WorkerSpec  `json:"workers"`
	Router  *RouterSpec  `json:"router"`
	Planner *PlannerSpec `json:"planner"`
}

// DynamoConfig defines the configuration for disaggregated dynamo workers
type DynamoConfig struct {
	// Name of the ConfigMap containing the worker config.
	Name string `json:"name"`
	// MountPath is the path where the config file should be mounted in the container.
	MountPath string `json:"mountPath"`
}

// WorkerSpec defines the spec for the deployment of dynamo prefill and decode workers
type WorkerSpec struct {
	Prefill *PrefillSpec `json:"prefill"`
	Decode  *DecoderSpec `json:"decode"`
}

// PrefillSpec defines the spec for the deployment of dynamo prefill workers
type PrefillSpec struct {
	Args []string        `json:"args,omitempty"`
	Env  []corev1.EnvVar `json:"env,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default:=1
	Replicas int `json:"replicas,omitempty"`
	// Resources is the resource requirements for the prefill worker.
	Resources *ResourceRequirements `json:"resources,omitempty"`
	// DRAResources is the list of DRA resource claims to be used for the prefill worker.
	DRAResources []DRAResource       `json:"draResources,omitempty"`
	Annotations  map[string]string   `json:"annotations,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
	PodAffinity  *corev1.PodAffinity `json:"podAffinity,omitempty"`
}

// DecoderSpec defines the spec for the deployment of dynamo decode workers
type DecoderSpec struct {
	Args []string        `json:"args,omitempty"`
	Env  []corev1.EnvVar `json:"env,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default:=1
	Replicas int `json:"replicas,omitempty"`
	// Resources is the resource requirements for the decode worker.
	Resources *ResourceRequirements `json:"resources,omitempty"`
	// DRAResources is the list of DRA resource claims to be used for the decode worker.
	DRAResources []DRAResource       `json:"draResources,omitempty"`
	Annotations  map[string]string   `json:"annotations,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
	PodAffinity  *corev1.PodAffinity `json:"podAffinity,omitempty"`
}

// RouterSpec defines the spec for the deployment of dynamo KV aware router
type RouterSpec struct {
	Args []string        `json:"args,omitempty"`
	Env  []corev1.EnvVar `json:"env,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default:=1
	Replicas     int                 `json:"replicas,omitempty"`
	Annotations  map[string]string   `json:"annotations,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
	PodAffinity  *corev1.PodAffinity `json:"podAffinity,omitempty"`
}

// PlannerSpec defines the spec for the deployment of dynamo planner
type PlannerSpec struct {
	Args []string        `json:"args,omitempty"`
	Env  []corev1.EnvVar `json:"env,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default:=1
	Replicas     int                 `json:"replicas,omitempty"`
	ScalePrefill *PrefillScaleSpec   `json:"prefillScale,omitempty"`
	ScaleDecode  *DecoderScaleSpec   `json:"decodeScale,omitempty"`
	Annotations  map[string]string   `json:"annotations,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
	PodAffinity  *corev1.PodAffinity `json:"podAffinity,omitempty"`
}

// PrefillScaleSpec defines the spec for autoscaling prefill workers
type PrefillScaleSpec struct {
	// TODO: define spec to configure replicas for prefill workers using KEDA
	HPA         HorizontalPodAutoscalerSpec `json:"hpa,omitempty"`
	Annotations map[string]string           `json:"annotations,omitempty"`
}

// DecoderScaleSpec defines the spec for autoscaling decode workers
type DecoderScaleSpec struct {
	// TODO: define spec to configure replicas for decode workers using KEDA
	HPA         HorizontalPodAutoscalerSpec `json:"hpa,omitempty"`
	Annotations map[string]string           `json:"annotations,omitempty"`
}

// DistributedNIMServiceStatus defines the observed state of DistributedNIMService
type DistributedNIMServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DistributedNIMService is the Schema for the distributednimservices API
type DistributedNIMService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DistributedNIMServiceSpec   `json:"spec,omitempty"`
	Status DistributedNIMServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DistributedNIMServiceList contains a list of DistributedNIMService
type DistributedNIMServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DistributedNIMService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DistributedNIMService{}, &DistributedNIMServiceList{})
}
