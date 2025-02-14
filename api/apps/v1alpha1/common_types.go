/*
Copyright 2024.

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
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// DefaultAPIPort is the default API  port
	DefaultAPIPort = 8000
	// DefaultNamedPortAPI is the default named API port
	DefaultNamedPortAPI = "api"
	// DefaultNamedPortMetrics is the default named Metrics port
	DefaultNamedPortMetrics = "metrics"
)

// Expose defines attributes to expose the service
type Expose struct {
	Service Service `json:"service"`
	Ingress Ingress `json:"ingress,omitempty"`
}

// Service defines attributes to create a service
type Service struct {
	Type corev1.ServiceType `json:"type,omitempty"`
	// override the default service name
	Name string `json:"name,omitempty"`
	// Deprecated: Use Ports instead.
	// +kubebuilder:deprecatedversion
	Port int32 `json:"port,omitempty"`
	// Defines multiple ports for the service
	Ports       []ServicePort     `json:"ports,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ServicePort defines attributes to setup the service ports
type ServicePort struct {
	// The name of this port within the service.
	Name string `json:"name,omitempty"`

	// The IP protocol for this port. Supports "TCP", "UDP", and "SCTP".
	// Default is TCP.
	// +kubebuilder:validation:Enum=TCP;UDP;SCTP
	// +kubebuilder:default="TCP"
	Protocol corev1.Protocol `json:"protocol,omitempty"`

	// The port that will be exposed by this service.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`
}

// Metrics defines attributes to setup metrics collection
type Metrics struct {
	Enabled *bool `json:"enabled,omitempty"`
	// for use with the Prometheus Operator and the primary service object
	ServiceMonitor ServiceMonitor `json:"serviceMonitor,omitempty"`
}

// ServiceMonitor defines attributes to create a service monitor
type ServiceMonitor struct {
	AdditionalLabels map[string]string `json:"additionalLabels,omitempty"`
	Annotations      map[string]string `json:"annotations,omitempty"`
	Interval         promv1.Duration   `json:"interval,omitempty"`
	ScrapeTimeout    promv1.Duration   `json:"scrapeTimeout,omitempty"`
}

// Autoscaling defines attributes to automatically scale the service based on metrics
type Autoscaling struct {
	Enabled     *bool                       `json:"enabled,omitempty"`
	HPA         HorizontalPodAutoscalerSpec `json:"hpa,omitempty"`
	Annotations map[string]string           `json:"annotations,omitempty"`
}

// HorizontalPodAutoscalerSpec defines the parameters required to setup HPA
type HorizontalPodAutoscalerSpec struct {
	MinReplicas *int32                                         `json:"minReplicas,omitempty"`
	MaxReplicas int32                                          `json:"maxReplicas"`
	Metrics     []autoscalingv2.MetricSpec                     `json:"metrics,omitempty"`
	Behavior    *autoscalingv2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty" `
}

// Image defines image attributes
type Image struct {
	Repository  string   `json:"repository"`
	PullPolicy  string   `json:"pullPolicy,omitempty"`
	Tag         string   `json:"tag"`
	PullSecrets []string `json:"pullSecrets,omitempty"`
}

// Ingress defines attributes to enable ingress for the service
type Ingress struct {
	// ingress, or virtualService - not both
	Enabled     *bool                    `json:"enabled,omitempty"`
	Annotations map[string]string        `json:"annotations,omitempty"`
	Spec        networkingv1.IngressSpec `json:"spec,omitempty"`
}

// IngressHost defines attributes for ingress host
type IngressHost struct {
	Host  string        `json:"host,omitempty"`
	Paths []IngressPath `json:"paths,omitempty"`
}

// IngressPath defines attributes for ingress paths
type IngressPath struct {
	Path        string                `json:"path,omitempty"`
	PathType    networkingv1.PathType `json:"pathType,omitempty"`
	ServiceType string                `json:"serviceType,omitempty"`
}

// Probe defines attributes for startup/liveness/readiness probes
type Probe struct {
	Enabled *bool         `json:"enabled,omitempty"`
	Probe   *corev1.Probe `json:"probe,omitempty"`
}

// CertConfig defines the configuration for custom certificates.
type CertConfig struct {
	// Name of the ConfigMap containing the certificate data.
	Name string `json:"name"`
	// MountPath is the path where the certificates should be mounted in the container.
	MountPath string `json:"mountPath"`
}

// selectNamedPort returns the first occurrence of a given named port, or an empty string if not found.
func selectNamedPort(serviceSpec Service, portNames ...string) string {
	for _, name := range portNames {
		for _, port := range serviceSpec.Ports {
			if port.Name == name {
				return name
			}
		}
	}
	return ""
}

// getProbePort determines the appropriate port for probes based on the service spec.
func getProbePort(serviceSpec Service) intstr.IntOrString {
	switch len(serviceSpec.Ports) {
	case 1:
		port := serviceSpec.Ports[0]
		if port.Name != "" {
			return intstr.FromString(port.Name)
		}
		return intstr.FromInt(int(port.Port))
	case 0:
		// Default to "api" as the operator always adds a default named port with 8000
		return intstr.FromString(DefaultNamedPortAPI)
	default:
		// Multiple ports: Prefer "api"
		if portName := selectNamedPort(serviceSpec, DefaultNamedPortAPI); portName != "" {
			return intstr.FromString(portName)
		}
		// Default when multiple ports exist
		return intstr.FromString(DefaultNamedPortAPI)
	}
}

// getMetricsPort determines the appropriate port for metrics based on the service spec.
func getMetricsPort(serviceSpec Service) intstr.IntOrString {
	switch len(serviceSpec.Ports) {
	case 1:
		port := serviceSpec.Ports[0]
		if port.Name != "" {
			return intstr.FromString(port.Name)
		}
		return intstr.FromInt(int(port.Port))
	case 0:
		// Default to "api" as the operator always adds a default named port with 8000
		return intstr.FromString(DefaultNamedPortAPI)
	default:
		// Multiple ports: Prefer "metrics", fallback to "api"
		if portName := selectNamedPort(serviceSpec, DefaultNamedPortMetrics, DefaultNamedPortAPI); portName != "" {
			return intstr.FromString(portName)
		}
		// Default when multiple ports exist
		return intstr.FromString(DefaultNamedPortMetrics)
	}
}
