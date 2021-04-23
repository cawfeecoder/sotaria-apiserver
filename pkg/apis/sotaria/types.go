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

package sotaria

import (
	v1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var swaggerMetadataDescriptionsObj = metav1.ObjectMeta{}.SwaggerDoc()

var swaggerMetadataDescriptionsNS = v1.NamespaceStatus{}.SwaggerDoc()

var (
	AdditionalPrinterColumns = []apiextensions.CustomResourceColumnDefinition{
		{
			Name:        "Status",
			Type:        "string",
			Description: swaggerMetadataDescriptionsNS["phase"],
			JSONPath:    ".status.phase",
		},
		{
			Name:        "Age",
			Type:        "string",
			Description: swaggerMetadataDescriptionsObj["creationTimestamp"],
			JSONPath:    ".objectMeta.creationTimestamp",
		},
	}
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FlunderList is a list of Flunder objects.
type ProjectList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Project
}

// FlunderSpec is the specification of a Flunder.
type ProjectSpec struct {
	Finalizers []v1.FinalizerName
}

// FlunderStatus is the status of a Flunder.
type ProjectStatus struct {
	Phase v1.NamespacePhase
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Flunder is an example type with a spec and a status.
type Project struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   ProjectSpec
	Status ProjectStatus
}
