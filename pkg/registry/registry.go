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

package registry

import (
	"context"
	"fmt"

	sotariaapi "github.com/nfrush/sotaria-apiserver/pkg/apis/sotaria"
	"github.com/nfrush/sotaria-apiserver/pkg/util"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	rbacv1listers "k8s.io/client-go/listers/rbac/v1"
)

// REST implements a RESTStorage for API services against etcd
type REST struct {
	*kubernetes.Clientset
	corev1listers.NamespaceLister
	rbacv1listers.ClusterRoleBindingLister
	*genericregistry.Store
}

// RESTInPeace is just a simple function that panics on error.
// Otherwise returns the given storage object. It is meant to be
// a wrapper for sotaria registries.
func RESTInPeace(storage rest.StandardStorage, err error) rest.StandardStorage {
	if err != nil {
		err = fmt.Errorf("unable to create REST storage for a resource due to %v, will die", err)
		panic(err)
	}
	return storage
}

func (s *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return s.Store.Get(ctx, name, options)
}

func (s *REST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		return nil, kerrors.NewForbidden(sotariaapi.Resource("project"), "", fmt.Errorf("unable to list projects without a user on the context"))
	}
	fmt.Printf("User: %v\n", user)
	projectSelector, err := labels.Parse("security.sotaria.io/type=project")
	if err != nil {
		return nil, err
	}
	clusteRoleBinding, err := s.ClusterRoleBindingLister.List(projectSelector)
	if err != nil {
		return nil, err
	}
	userClusterRoles := []string{}
	for _, v := range clusteRoleBinding {
		for _, z := range v.Subjects {
			fmt.Printf("Subjects: %z", z)
			if (z.Kind == "User" && z.Name == user.GetName()) || (z.Kind == "Group" && util.StrArrayContains(user.GetGroups(), z.Name)) || (z.Kind == "ServiceAccount" && z.Name == user.GetName()) {
				userClusterRoles = append(userClusterRoles, v.RoleRef.Name)
			}
		}
	}
	namespaces, err := s.NamespaceLister.List(projectSelector)
	if err != nil {
		return nil, err
	}
	projects := &sotariaapi.ProjectList{}
	isAdmin := util.StrArrayContains(user.GetGroups(), "system:masters")
	for _, namespace := range namespaces {
		if util.StrArrayContains(user.GetGroups(), namespace.Labels["security.sotaria.io/role"]) || util.StrArrayContains(userClusterRoles, namespace.Labels["security.sotaria.io/role"]) || isAdmin {
			project := sotariaapi.Project{
				ObjectMeta: namespace.ObjectMeta,
				Spec: sotariaapi.ProjectSpec{
					Finalizers: namespace.Spec.Finalizers,
				},
				Status: sotariaapi.ProjectStatus{
					Phase: namespace.Status.Phase,
				},
			}
			projects.Items = append(projects.Items, project)
		}
	}
	return projects, nil
}

func (s *REST) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	return s.Store.Watch(ctx, options)
}

func (s *REST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return s.Store.Update(ctx, name, objInfo, createValidation, updateValidation, forceAllowCreate, options)
}

func (s *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	return s.Store.Delete(ctx, name, deleteValidation, options)
}

func (s *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	return s.Store.Create(ctx, obj, createValidation, options)
}
