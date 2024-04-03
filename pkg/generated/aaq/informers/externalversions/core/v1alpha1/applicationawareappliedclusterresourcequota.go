/*
Copyright 2023 The AAQ Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	versioned "kubevirt.io/application-aware-quota/pkg/generated/aaq/clientset/versioned"
	internalinterfaces "kubevirt.io/application-aware-quota/pkg/generated/aaq/informers/externalversions/internalinterfaces"
	v1alpha1 "kubevirt.io/application-aware-quota/pkg/generated/aaq/listers/core/v1alpha1"
	corev1alpha1 "kubevirt.io/application-aware-quota/staging/src/kubevirt.io/application-aware-quota-api/pkg/apis/core/v1alpha1"
)

// ApplicationAwareAppliedClusterResourceQuotaInformer provides access to a shared informer and lister for
// ApplicationAwareAppliedClusterResourceQuotas.
type ApplicationAwareAppliedClusterResourceQuotaInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ApplicationAwareAppliedClusterResourceQuotaLister
}

type applicationAwareAppliedClusterResourceQuotaInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewApplicationAwareAppliedClusterResourceQuotaInformer constructs a new informer for ApplicationAwareAppliedClusterResourceQuota type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewApplicationAwareAppliedClusterResourceQuotaInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredApplicationAwareAppliedClusterResourceQuotaInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredApplicationAwareAppliedClusterResourceQuotaInformer constructs a new informer for ApplicationAwareAppliedClusterResourceQuota type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredApplicationAwareAppliedClusterResourceQuotaInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AaqV1alpha1().ApplicationAwareAppliedClusterResourceQuotas(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AaqV1alpha1().ApplicationAwareAppliedClusterResourceQuotas(namespace).Watch(context.TODO(), options)
			},
		},
		&corev1alpha1.ApplicationAwareAppliedClusterResourceQuota{},
		resyncPeriod,
		indexers,
	)
}

func (f *applicationAwareAppliedClusterResourceQuotaInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredApplicationAwareAppliedClusterResourceQuotaInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *applicationAwareAppliedClusterResourceQuotaInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&corev1alpha1.ApplicationAwareAppliedClusterResourceQuota{}, f.defaultInformer)
}

func (f *applicationAwareAppliedClusterResourceQuotaInformer) Lister() v1alpha1.ApplicationAwareAppliedClusterResourceQuotaLister {
	return v1alpha1.NewApplicationAwareAppliedClusterResourceQuotaLister(f.Informer().GetIndexer())
}