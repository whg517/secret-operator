/*
Copyright 2024 zncdatadev.

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
)

// SecretClassSpec defines the desired state of SecretClass
type SecretClassSpec struct {
	Backend *BackendSpec `json:"backend,omitempty"`
}

type BackendSpec struct {
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	AutoTls *AutoTlsSpec `json:"autoTls,omitempty"`
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	K8sSearch *K8sSearchSpec `json:"k8sSearch,omitempty"`
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Kerberos *KerberosSpec `json:"kerberos,omitempty"`
}

type AutoTlsSpec struct {
	CA *CASpec `json:"ca,omitempty"`

	// Use time.ParseDuration to parse the string
	// Default is 360h (15 days)
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="360h"
	MaxCertificateLifeTime string `json:"maxCertificateLifeTime,omitempty"`
}

type CASpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	AutoGenerated bool `json:"autoGenerated,omitempty"`

	// Use time.ParseDuration to parse the string
	// Default is 8760h (1 year)
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="8760h"
	CACertificateLifeTime string `json:"caCertificateLifeTime,omitempty"`

	// +kubebuilder:validation:Required
	Secret *SecretSpec `json:"secret,omitempty"`
}

type SecretSpec struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// TODO implement the KerberosSpec
type KerberosSpec struct {
}

type K8sSearchSpec struct {
	// +kubebuilder:validation:Required
	SearchNamespace *SearchNamespaceSpec `json:"searchNamespace,omitempty"`
}

type SearchNamespaceSpec struct {
	Name *string `json:"name,omitempty"`

	Pod *PodSpec `json:"pod,omitempty"`
}

type PodSpec struct {
}

// SecretClassStatus defines the observed state of SecretClass
type SecretClassStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=secretclasses,scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +operator-sdk:csv:customresourcedefinitions:displayName="Secret Class"

// SecretClass is the Schema for the secretclasses API
type SecretClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretClassSpec   `json:"spec,omitempty"`
	Status SecretClassStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SecretClassList contains a list of SecretClass
type SecretClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretClass{}, &SecretClassList{})
}
