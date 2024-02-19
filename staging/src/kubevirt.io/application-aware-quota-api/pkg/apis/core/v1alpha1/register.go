/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 *
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	applicationAwareResourceQuota "kubevirt.io/application-aware-quota/staging/src/kubevirt.io/application-aware-quota-api/pkg/apis/core"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion                                   = schema.GroupVersion{Group: applicationAwareResourceQuota.GroupName, Version: applicationAwareResourceQuota.LatestVersion}
	ApplicationAwareResourceQuotaGroupVersionKind        = schema.GroupVersionKind{Group: applicationAwareResourceQuota.GroupName, Version: applicationAwareResourceQuota.LatestVersion, Kind: "ApplicationAwareResourceQuota"}
	ApplicationAwareClusterResourceQuotaGroupVersionKind = schema.GroupVersionKind{Group: applicationAwareResourceQuota.GroupName, Version: applicationAwareResourceQuota.LatestVersion, Kind: "ApplicationAwareClusterResourceQuota"}
)

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder initializes a scheme builder
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ApplicationAwareResourceQuota{},
		&ApplicationAwareResourceQuotaList{},
		&ApplicationAwareClusterResourceQuota{},
		&ApplicationAwareClusterResourceQuotaList{},
		&ApplicationAwareAppliedClusterResourceQuota{},
		&ApplicationAwareAppliedClusterResourceQuotaList{},
		&AAQJobQueueConfig{},
		&AAQJobQueueConfigList{},
		&AAQ{},
		&AAQList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func AddKnownTypesGenerator(groupVersions []schema.GroupVersion) func(scheme *runtime.Scheme) error {
	// Adds the list of known types to api.Scheme.
	return func(scheme *runtime.Scheme) error {
		for _, groupVersion := range groupVersions {
			scheme.AddKnownTypes(groupVersion,
				&ApplicationAwareResourceQuota{},
				&ApplicationAwareResourceQuotaList{},
				&ApplicationAwareClusterResourceQuota{},
				&ApplicationAwareClusterResourceQuotaList{},
				&ApplicationAwareAppliedClusterResourceQuota{},
				&ApplicationAwareAppliedClusterResourceQuotaList{},
				&AAQJobQueueConfig{},
				&AAQJobQueueConfigList{},
				&AAQ{},
				&AAQList{},
			)
			metav1.AddToGroupVersion(scheme, groupVersion)
		}

		return nil
	}
}
