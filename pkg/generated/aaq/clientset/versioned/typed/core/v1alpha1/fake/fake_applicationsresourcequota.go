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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "kubevirt.io/applications-aware-quota/staging/src/kubevirt.io/applications-aware-quota-api/pkg/apis/core/v1alpha1"
)

// FakeApplicationsResourceQuotas implements ApplicationsResourceQuotaInterface
type FakeApplicationsResourceQuotas struct {
	Fake *FakeAaqV1alpha1
	ns   string
}

var applicationsresourcequotasResource = schema.GroupVersionResource{Group: "aaq.kubevirt.io", Version: "v1alpha1", Resource: "applicationsresourcequotas"}

var applicationsresourcequotasKind = schema.GroupVersionKind{Group: "aaq.kubevirt.io", Version: "v1alpha1", Kind: "ApplicationsResourceQuota"}

// Get takes name of the applicationsResourceQuota, and returns the corresponding applicationsResourceQuota object, and an error if there is any.
func (c *FakeApplicationsResourceQuotas) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ApplicationsResourceQuota, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(applicationsresourcequotasResource, c.ns, name), &v1alpha1.ApplicationsResourceQuota{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationsResourceQuota), err
}

// List takes label and field selectors, and returns the list of ApplicationsResourceQuotas that match those selectors.
func (c *FakeApplicationsResourceQuotas) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ApplicationsResourceQuotaList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(applicationsresourcequotasResource, applicationsresourcequotasKind, c.ns, opts), &v1alpha1.ApplicationsResourceQuotaList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ApplicationsResourceQuotaList{ListMeta: obj.(*v1alpha1.ApplicationsResourceQuotaList).ListMeta}
	for _, item := range obj.(*v1alpha1.ApplicationsResourceQuotaList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested applicationsResourceQuotas.
func (c *FakeApplicationsResourceQuotas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(applicationsresourcequotasResource, c.ns, opts))

}

// Create takes the representation of a applicationsResourceQuota and creates it.  Returns the server's representation of the applicationsResourceQuota, and an error, if there is any.
func (c *FakeApplicationsResourceQuotas) Create(ctx context.Context, applicationsResourceQuota *v1alpha1.ApplicationsResourceQuota, opts v1.CreateOptions) (result *v1alpha1.ApplicationsResourceQuota, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(applicationsresourcequotasResource, c.ns, applicationsResourceQuota), &v1alpha1.ApplicationsResourceQuota{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationsResourceQuota), err
}

// Update takes the representation of a applicationsResourceQuota and updates it. Returns the server's representation of the applicationsResourceQuota, and an error, if there is any.
func (c *FakeApplicationsResourceQuotas) Update(ctx context.Context, applicationsResourceQuota *v1alpha1.ApplicationsResourceQuota, opts v1.UpdateOptions) (result *v1alpha1.ApplicationsResourceQuota, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(applicationsresourcequotasResource, c.ns, applicationsResourceQuota), &v1alpha1.ApplicationsResourceQuota{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationsResourceQuota), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeApplicationsResourceQuotas) UpdateStatus(ctx context.Context, applicationsResourceQuota *v1alpha1.ApplicationsResourceQuota, opts v1.UpdateOptions) (*v1alpha1.ApplicationsResourceQuota, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(applicationsresourcequotasResource, "status", c.ns, applicationsResourceQuota), &v1alpha1.ApplicationsResourceQuota{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationsResourceQuota), err
}

// Delete takes name of the applicationsResourceQuota and deletes it. Returns an error if one occurs.
func (c *FakeApplicationsResourceQuotas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(applicationsresourcequotasResource, c.ns, name, opts), &v1alpha1.ApplicationsResourceQuota{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApplicationsResourceQuotas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(applicationsresourcequotasResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ApplicationsResourceQuotaList{})
	return err
}

// Patch applies the patch and returns the patched applicationsResourceQuota.
func (c *FakeApplicationsResourceQuotas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ApplicationsResourceQuota, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(applicationsresourcequotasResource, c.ns, name, pt, data, subresources...), &v1alpha1.ApplicationsResourceQuota{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationsResourceQuota), err
}