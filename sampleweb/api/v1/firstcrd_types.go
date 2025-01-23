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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FirstCrdSpec defines the desired state of FirstCrd.
type FirstCrdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of FirstCrd. Edit firstcrd_types.go to remove/update
	Foo string `json:"foo,omitempty"`

	Message  string `json:"message"`  // Webページに出力するメッセージ
	Replicas int32  `json:"replicas"` // Deploymentのレプリカ数
}

// FirstCrdStatus defines the observed state of FirstCrd.
type FirstCrdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	CurrentMessage  string `json:"message"`         // Webページに出力するメッセージ
	CurrentReplicas int32  `json:"currentReplicas"` // Deploymentのレプリカ数
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// FirstCrd is the Schema for the firstcrds API.
type FirstCrd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FirstCrdSpec   `json:"spec,omitempty"`
	Status FirstCrdStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FirstCrdList contains a list of FirstCrd.
type FirstCrdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FirstCrd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FirstCrd{}, &FirstCrdList{})
}
