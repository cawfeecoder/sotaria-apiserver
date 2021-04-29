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

package project

import (
	"github.com/nfrush/sotaria-apiserver/pkg/apis/sotaria"
	"github.com/nfrush/sotaria-apiserver/pkg/registry"
	"github.com/nfrush/sotaria-apiserver/pkg/registry/tableconvertor"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	rbacv1listers "k8s.io/client-go/listers/rbac/v1"
)

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter, nsLister corev1listers.NamespaceLister, crbLister rbacv1listers.ClusterRoleBindingLister, kubeClient *kubernetes.Clientset) (*registry.REST, error) {
	strategy := NewStrategy(scheme)

	table_convert, err := tableconvertor.New(sotaria.AdditionalPrinterColumns)

	if err != nil {
		return nil, err
	}

	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &sotaria.Project{} },
		NewListFunc:              func() runtime.Object { return &sotaria.ProjectList{} },
		PredicateFunc:            MatchProject,
		DefaultQualifiedResource: sotaria.Resource("projects"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: table_convert,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &registry.REST{kubeClient, nsLister, crbLister, store}, nil
}
