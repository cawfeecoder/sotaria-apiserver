/*
Copyright 2017 The Kubernetes Authors.

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProjectList is a list of Project objects.
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Project `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ProjectSpec struct {
	Finalizers []v1.FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=k8s.io/kubernetes/pkg/api/v1.FinalizerName"`
}

type ProjectStatus struct {
	Phase v1.NamespacePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=k8s.io/kubernetes/pkg/api/v1.NamespacePhase"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ProjectSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status ProjectStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
